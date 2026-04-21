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
	Datasets           []Dataset  `json:"datasets,omitempty"`
}

// Dataset holds filesystem, volume, or snapshot metrics within a pool.
type Dataset struct {
	Name              string  `json:"name"`
	Pool              string  `json:"pool"`
	Type              string  `json:"type"`
	Available         float64 `json:"available"`
	LogicalUsed       float64 `json:"logical_used"`
	Quota             float64 `json:"quota"`
	Referenced        float64 `json:"referenced"`
	UsedByDataset     float64 `json:"used_by_dataset"`
	Used              float64 `json:"used"`
	VolumeSize        float64 `json:"volume_size"`
	Written           float64 `json:"written"`
	UsedPercent       float64 `json:"used_percent"`
	QuotaUsedPercent  float64 `json:"quota_used_percent"`
	VolumeUsedPercent float64 `json:"volume_used_percent"`
}

// NodeData holds all pool data fetched from one endpoint.
type NodeData struct {
	Label     string    `json:"label"`
	Location  string    `json:"location,omitempty"`
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
	datasets := map[string]*Dataset{}
	ensure := func(name string) *Pool {
		if _, ok := pools[name]; !ok {
			pools[name] = &Pool{Name: name, Health: HealthUnknown}
		}
		return pools[name]
	}
	ensureDataset := func(pool, name, datasetType string) *Dataset {
		key := pool + "\x00" + name
		if _, ok := datasets[key]; !ok {
			datasets[key] = &Dataset{
				Name: name,
				Pool: pool,
				Type: datasetType,
			}
		}
		if datasetType != "" {
			datasets[key].Type = datasetType
		}
		return datasets[key]
	}

	for _, s := range samples {
		pool := s.Labels["pool"]
		if pool == "" {
			continue
		}
		p := ensure(pool)
		if datasetName := s.Labels["name"]; datasetName != "" {
			d := ensureDataset(pool, datasetName, s.Labels["type"])
			switch s.Name {
			case "zfs_dataset_available_bytes":
				d.Available = s.Value
			case "zfs_dataset_logical_used_bytes":
				d.LogicalUsed = s.Value
			case "zfs_dataset_quota_bytes":
				d.Quota = s.Value
			case "zfs_dataset_referenced_bytes":
				d.Referenced = s.Value
			case "zfs_dataset_used_by_dataset_bytes":
				d.UsedByDataset = s.Value
			case "zfs_dataset_used_bytes":
				d.Used = s.Value
			case "zfs_dataset_volume_size_bytes":
				d.VolumeSize = s.Value
			case "zfs_dataset_written_bytes":
				d.Written = s.Value
			}
		}

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
	for _, d := range datasets {
		if d.Used > 0 && d.Available > 0 {
			total := d.Used + d.Available
			d.UsedPercent = math.Round((d.Used/total)*10000) / 100
		}
		if d.Quota > 0 {
			d.QuotaUsedPercent = math.Round((d.Used/d.Quota)*10000) / 100
		}
		if d.VolumeSize > 0 {
			d.VolumeUsedPercent = math.Round((d.Used/d.VolumeSize)*10000) / 100
		}
		if p, ok := pools[d.Pool]; ok {
			p.Datasets = append(p.Datasets, *d)
		}
	}

	for _, p := range pools {
		if p.Size > 0 {
			if p.Allocated > 0 && p.Free == 0 {
				p.Free = p.Size - p.Allocated
			} else if p.Free > 0 && p.Allocated == 0 {
				p.Allocated = p.Size - p.Free
			}
			p.UsedPercent = math.Round((p.Allocated/p.Size)*10000) / 100
		}
		slices.SortFunc(p.Datasets, func(a, b Dataset) int {
			if a.Used != b.Used {
				if a.Used > b.Used {
					return -1
				}
				return 1
			}
			return strings.Compare(a.Name, b.Name)
		})
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
