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
func NewRecorder(store *Store, f *fetcher.Fetcher, interval time.Duration) *Recorder {
	return &Recorder{store: store, fetcher: f, interval: interval}
}

// Start begins the background recording loop; returns immediately.
func (r *Recorder) Start(ctx context.Context) {
	go r.run(ctx)
}

func (r *Recorder) run(ctx context.Context) {
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
					Key: SeriesKey(node.Label, "pool", pool.Name, "alloc_bytes"),
					Ts:  now, Value: pool.Allocated,
				})
			}
			if pool.Free > 0 {
				samples = append(samples, Sample{
					Key: SeriesKey(node.Label, "pool", pool.Name, "free_bytes"),
					Ts:  now, Value: pool.Free,
				})
			}
		}
		for _, disk := range node.Disks {
			if disk.Temperature > 0 {
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
