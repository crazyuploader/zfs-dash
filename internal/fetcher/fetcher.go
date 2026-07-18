// Package fetcher fetches Prometheus metrics from multiple endpoints concurrently.
package fetcher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/crazyuploader/zfs-dash/internal/config"
	"github.com/crazyuploader/zfs-dash/internal/model"
	"github.com/crazyuploader/zfs-dash/internal/parser"
)

const fetchTimeout = 10 * time.Second
const maxResponseBytes = 10 << 20

// Fetcher retrieves metrics from configured endpoints.
type Fetcher struct {
	client    *http.Client
	mu        sync.RWMutex
	endpoints []config.Endpoint
	gen       uint64 // bumped by SetEndpoints; stale fetches are discarded
	cacheTTL  time.Duration
	cache     []model.NodeData
	expiresAt time.Time
	rates     rateTracker
}

// New creates a Fetcher for the provided endpoints.
func New(endpoints []config.Endpoint, cacheTTL time.Duration) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: fetchTimeout,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		endpoints: endpoints,
		cacheTTL:  cacheTTL,
	}
}

// SetEndpoints updates the target list (for hot-reload).
func (f *Fetcher) SetEndpoints(eps []config.Endpoint) {
	urls := make(map[string]struct{}, len(eps))
	for _, ep := range eps {
		if ep.NodeExporterURL != "" {
			urls[ep.NodeExporterURL] = struct{}{}
		}
	}
	f.rates.retain(urls)

	f.mu.Lock()
	defer f.mu.Unlock()
	f.endpoints = eps
	f.gen++
	f.expiresAt = time.Time{} // invalidates cache
}

// CacheInfo returns the current cache status.
func (f *Fetcher) CacheInfo() (expiresAt time.Time, ttl time.Duration) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.expiresAt, f.cacheTTL
}

// FetchAll fetches all endpoints concurrently and returns results in the same order.
// It returns the results and a boolean indicating if the results were from cache.
// Results share inner slices with the cache and must be treated as read-only.
func (f *Fetcher) FetchAll(ctx context.Context) ([]model.NodeData, bool) {
	f.mu.RLock()
	if time.Now().Before(f.expiresAt) {
		slog.Debug("cache HIT", "expires_in", time.Until(f.expiresAt).Round(time.Second))
		// Shallow copy: inner slices (Pools, Disks, Datasets) are shared with
		// the cache — callers must treat results as read-only.
		data := append([]model.NodeData{}, f.cache...)
		f.mu.RUnlock()
		return data, true
	}
	f.mu.RUnlock()

	return f.fetchAndStore(ctx), false
}

// Refresh fetches all endpoints unconditionally (bypassing the cache TTL)
// and updates the cache. Used by the background poller so its cadence is
// independent of cache_ttl. Results are read-only (see FetchAll).
func (f *Fetcher) Refresh(ctx context.Context) []model.NodeData {
	return f.fetchAndStore(ctx)
}

// fetchAndStore runs the concurrent fan-out WITHOUT holding f.mu during network
// I/O: it snapshots the endpoint list and generation under a read lock,
// fetches unlocked, then publishes the results only if the endpoints have
// not been swapped by SetEndpoints in the meantime.
func (f *Fetcher) fetchAndStore(ctx context.Context) []model.NodeData {
	f.mu.RLock()
	eps := append([]config.Endpoint{}, f.endpoints...)
	gen := f.gen
	f.mu.RUnlock()

	slog.Debug("fetching all endpoints", "endpoints", len(eps))
	results := make([]model.NodeData, len(eps))
	var wg sync.WaitGroup
	for i, ep := range eps {
		wg.Add(1)
		go func(i int, ep config.Endpoint) {
			defer wg.Done()
			results[i] = f.fetchOne(ctx, ep)
		}(i, ep)
	}
	wg.Wait()

	f.mu.Lock()
	if gen == f.gen {
		f.cache = results
		f.expiresAt = time.Now().Add(f.cacheTTL)
	}
	// Stale generation: endpoints changed mid-fetch; return the results to
	// this caller but do not overwrite the newer configuration's cache.
	f.mu.Unlock()

	// Shallow copy; results are read-only (see above).
	return append([]model.NodeData{}, results...)
}

// fetchRaw fetches and parses Prometheus text-format metrics from a single URL.
func (f *Fetcher) fetchRaw(ctx context.Context, label, url string) ([]parser.Sample, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	resp, err := f.client.Do(req)
	if err != nil {
		slog.Warn("fetch failed", "label", label, "url", url, "error", err)
		return nil, fmt.Errorf("unreachable: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("fetch failed", "label", label, "url", url, "status", resp.StatusCode)
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil {
		slog.Warn("fetch read error", "label", label, "url", url, "error", err)
		return nil, fmt.Errorf("read: %w", err)
	}
	if len(body) > maxResponseBytes {
		slog.Warn("fetch response too large", "label", label, "url", url, "limit", maxResponseBytes)
		return nil, fmt.Errorf("response too large: limit %d bytes", maxResponseBytes)
	}
	slog.Debug("read metrics", "label", label, "url", url, "bytes", len(body))
	samples, err := parser.Parse(bytes.NewReader(body))
	if err != nil {
		slog.Warn("parse failed", "label", label, "url", url, "error", err)
		return nil, fmt.Errorf("parse: %w", err)
	}
	return samples, nil
}

func (f *Fetcher) fetchOne(ctx context.Context, ep config.Endpoint) model.NodeData {
	slog.Debug("fetching metrics", "label", ep.Label, "url", ep.URL)
	nd := model.NodeData{
		Label:     ep.Label,
		Location:  ep.Location,
		URL:       ep.URL,
		FetchedAt: time.Now(),
	}

	var (
		zfsSamples      []parser.Sample
		smartctlSamples []parser.Sample
		nodeSamples     []parser.Sample
		zfsErr          error
	)

	// Fetch ZFS plus optional companion exporters concurrently;
	// companion failures are non-fatal.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		zfsSamples, zfsErr = f.fetchRaw(ctx, ep.Label, ep.URL)
	}()
	if ep.SmartctlURL != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			smartctlSamples, err = f.fetchRaw(ctx, ep.Label, ep.SmartctlURL)
			if err != nil {
				slog.Warn("smartctl fetch failed (disk data unavailable)", "label", ep.Label, "error", err)
			}
		}()
	}
	if ep.NodeExporterURL != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			nodeSamples, err = f.fetchRaw(ctx, ep.Label, ep.NodeExporterURL)
			if err != nil {
				slog.Warn("node_exporter fetch failed (system data unavailable)", "label", ep.Label, "error", err)
			}
		}()
	}
	wg.Wait()

	if len(nodeSamples) > 0 {
		nd.System = model.ExtractSystem(nodeSamples)
		f.rates.apply(ep.NodeExporterURL, nd.System)
	}

	if zfsErr != nil {
		nd.Error = zfsErr.Error()
		// Still populate disks from smartctl even when ZFS is unavailable.
		if len(smartctlSamples) > 0 {
			nd.Disks = model.ExtractDisks(smartctlSamples)
			nd.SmartctlInfo = model.ExtractSmartctlInfo(smartctlSamples)
		}
		return nd
	}

	allSamples := make([]parser.Sample, 0, len(zfsSamples)+len(smartctlSamples))
	allSamples = append(allSamples, zfsSamples...)
	allSamples = append(allSamples, smartctlSamples...)
	nd.Pools = model.ExtractPools(allSamples)
	nd.ExporterInfo = model.ExtractExporterInfo(allSamples)
	nd.Disks = model.ExtractDisks(allSamples)
	if len(smartctlSamples) > 0 {
		nd.SmartctlInfo = model.ExtractSmartctlInfo(smartctlSamples)
	}
	slog.Debug("extracted pools", "label", ep.Label, "count", len(nd.Pools), "disks", len(nd.Disks))
	return nd
}
