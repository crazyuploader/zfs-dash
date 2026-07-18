package history

import (
	"context"
	"log/slog"
	"time"

	"github.com/crazyuploader/zfs-dash/internal/fetcher"
)

// Recorder polls the fetcher and writes metrics to the Store.
type Recorder struct {
	store    *Store
	fetcher  *fetcher.Fetcher
	interval time.Duration
}

// NewRecorder creates a Recorder that polls at the given interval.
// interval is clamped to a minimum of 1 second to prevent a zero/negative
// duration from causing time.NewTicker to panic.
func NewRecorder(store *Store, f *fetcher.Fetcher, interval time.Duration) *Recorder {
	if interval < time.Second {
		interval = time.Second
	}
	return &Recorder{store: store, fetcher: f, interval: interval}
}

// Run begins the recording loop; blocks until ctx is cancelled.
// Callers that want a background loop should invoke it as "go rec.Run(ctx)".
func (r *Recorder) Run(ctx context.Context) {
	// Initial prune on startup to clean up stale data.
	if err := r.store.Prune(); err != nil {
		slog.Warn("history prune failed", "error", err)
	}

	// Record immediately, then on every tick.
	r.record(ctx)

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	pruneTicker := time.NewTicker(1 * time.Hour)
	defer pruneTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.record(ctx)
		case <-pruneTicker.C:
			if err := r.store.Prune(); err != nil {
				slog.Warn("history prune failed", "error", err)
			}
		}
	}
}

func (r *Recorder) record(ctx context.Context) {
	nodes, _ := r.fetcher.FetchAll(ctx)
	now := time.Now()
	var samples []Sample

	for _, node := range nodes {
		for _, pool := range node.Pools {
			samples = append(samples, Sample{
				Key: SeriesKey(node.Label, "pool", pool.Name, "used_pct"),
				Ts:  now, Value: pool.UsedPercent,
			})
			if pool.Allocated > 0 {
				samples = append(samples, Sample{
					Key:   SeriesKey(node.Label, "pool", pool.Name, "alloc_bytes"),
					Ts:    now,
					Value: pool.Allocated / (1 << 20), // Store in MiB to avoid float32 precision loss
				})
			}
			if pool.Free > 0 {
				samples = append(samples, Sample{
					Key:   SeriesKey(node.Label, "pool", pool.Name, "free_bytes"),
					Ts:    now,
					Value: pool.Free / (1 << 20), // Store in MiB to avoid float32 precision loss
				})
			}
		}
		if sys := node.System; sys != nil {
			// Fixed "node" name component keeps series keys stable even if
			// the reported hostname changes.
			if sys.HasCPURates {
				samples = append(samples, Sample{
					Key: SeriesKey(node.Label, "system", "node", "cpu_pct"),
					Ts:  now, Value: sys.CPUBusyPct,
				})
			}
			if sys.MemTotal > 0 {
				samples = append(samples, Sample{
					Key: SeriesKey(node.Label, "system", "node", "mem_used_pct"),
					Ts:  now, Value: sys.MemUsedPct,
				})
			}
			if sys.HasLoad {
				samples = append(samples, Sample{
					Key: SeriesKey(node.Label, "system", "node", "load1"),
					Ts:  now, Value: sys.Load1,
				})
			}
		}
		for _, disk := range node.Disks {
			if disk.HasTemperature {
				samples = append(samples, Sample{
					Key: SeriesKey(node.Label, "disk", disk.Device, "temp_c"),
					Ts:  now, Value: disk.Temperature,
				})
			}
			if disk.HasPercentUsed {
				samples = append(samples, Sample{
					Key: SeriesKey(node.Label, "disk", disk.Device, "wear_pct"),
					Ts:  now, Value: disk.PercentageUsed,
				})
			}
			if disk.HasWearLeveling {
				samples = append(samples, Sample{
					Key: SeriesKey(node.Label, "disk", disk.Device, "wear_lvl"),
					Ts:  now, Value: disk.WearLevelingCount,
				})
			}
			if disk.PowerOnHours > 0 {
				samples = append(samples, Sample{
					Key: SeriesKey(node.Label, "disk", disk.Device, "pow_hrs"),
					Ts:  now, Value: disk.PowerOnHours,
				})
			}
		}
	}

	if len(samples) == 0 {
		return
	}
	if err := r.store.WriteBatch(samples); err != nil {
		slog.Warn("history write failed", "error", err)
	} else {
		slog.Debug("history recorded", "samples", len(samples))
	}
}
