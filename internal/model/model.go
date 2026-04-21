// Package model defines ZFS domain types and metric extraction logic.
package model

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/crazyuploader/zfs-dash/internal/parser"
)

// PoolHealth is the human-readable pool state.
type PoolHealth string

const (
	HealthOnline    PoolHealth = "ONLINE"
	HealthDegraded  PoolHealth = "DEGRADED"
	HealthFaulted   PoolHealth = "FAULTED"
	HealthOffline   PoolHealth = "OFFLINE"
	HealthUnavail   PoolHealth = "UNAVAIL"
	HealthRemoved   PoolHealth = "REMOVED"
	HealthSuspended PoolHealth = "SUSPENDED"
	HealthUnknown   PoolHealth = "UNKNOWN"
)

// Pool holds the key metrics for one ZFS pool.
type Pool struct {
	Name               string     `json:"name"`
	Health             PoolHealth `json:"health"`
	Size               float64    `json:"size"`
	Allocated          float64    `json:"allocated"`
	Free               float64    `json:"free"`
	Freeing            float64    `json:"freeing"`
	LeakedBytes        float64    `json:"leaked_bytes"`
	DedupRatio         float64    `json:"dedup_ratio"`
	FragmentationRatio float64    `json:"fragmentation_ratio"`
	ReadOnly           bool       `json:"read_only"`
	UsedPercent        float64    `json:"used_percent"`
}

// NodeData holds all pool data fetched from one endpoint.
type NodeData struct {
	Label     string    `json:"label"`
	URL       string    `json:"url"`
	FetchedAt time.Time `json:"fetched_at"`
	Error     string    `json:"error,omitempty"`
	Pools     []Pool    `json:"pools"`
}

func healthFromValue(v float64) PoolHealth {
	switch int(math.Round(v)) {
	case 0:
		return HealthOnline
	case 1:
		return HealthDegraded
	case 2:
		return HealthFaulted
	case 3:
		return HealthOffline
	case 4:
		return HealthUnavail
	case 5:
		return HealthRemoved
	case 6:
		return HealthSuspended
	default:
		return HealthUnknown
	}
}

// ExtractPools builds Pool structs from a flat Prometheus sample slice.
func ExtractPools(samples []parser.Sample) []Pool {
	pools := map[string]*Pool{}
	ensure := func(name string) *Pool {
		if _, ok := pools[name]; !ok {
			pools[name] = &Pool{Name: name, Health: HealthUnknown}
		}
		return pools[name]
	}

	for _, s := range samples {
		pool := s.Labels["pool"]
		if pool == "" {
			continue
		}
		p := ensure(pool)
		switch s.Name {
		case "zfs_pool_health":
			p.Health = healthFromValue(s.Value)
		case "zfs_pool_size_bytes":
			p.Size = s.Value
		case "zfs_pool_allocated_bytes":
			p.Allocated = s.Value
		case "zfs_pool_free_bytes":
			p.Free = s.Value
		case "zfs_pool_freeing_bytes":
			p.Freeing = s.Value
		case "zfs_pool_leaked_bytes":
			p.LeakedBytes = s.Value
		case "zfs_pool_deduplication_ratio":
			p.DedupRatio = s.Value
		case "zfs_pool_fragmentation_ratio":
			p.FragmentationRatio = s.Value
		case "zfs_pool_readonly":
			p.ReadOnly = s.Value != 0
		}
	}

	result := make([]Pool, 0, len(pools))
	for _, p := range pools {
		if p.Size > 0 {
			if p.Allocated > 0 && p.Free == 0 {
				p.Free = p.Size - p.Allocated
			} else if p.Free > 0 && p.Allocated == 0 {
				p.Allocated = p.Size - p.Free
			}
			p.UsedPercent = math.Round((p.Allocated/p.Size)*10000) / 100
		}
		result = append(result, *p)
	}

	slices.SortFunc(result, func(a, b Pool) int {
		return strings.Compare(a.Name, b.Name)
	})

	return result
}

// HumanBytes returns a human-readable byte string (e.g. "3.72 TB").
func HumanBytes(b float64) string {
	const unit = 1024.0
	if b < unit {
		return fmt.Sprintf("%.0f B", b)
	}
	div, exp := unit, 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %s", b/div, []string{"KB", "MB", "GB", "TB", "PB", "EB"}[exp])
}
