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
	cacheTTL  time.Duration
	cache     []model.NodeData
	expiresAt time.Time
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
	f.mu.Lock()
	defer f.mu.Unlock()
	f.endpoints = eps
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
func (f *Fetcher) FetchAll(ctx context.Context) ([]model.NodeData, bool) {
	f.mu.RLock()
	if time.Now().Before(f.expiresAt) {
		slog.Debug("cache HIT", "expires_in", time.Until(f.expiresAt).Round(time.Second))
		// Return a copy so callers cannot mutate cached data.
		data := append([]model.NodeData{}, f.cache...)
		f.mu.RUnlock()
		return data, true
	}
	f.mu.RUnlock()

	// Single-flight like behavior could be added here, but for now just standard lock.
	f.mu.Lock()
	defer f.mu.Unlock()

	// Re-check after acquiring write lock
	if time.Now().Before(f.expiresAt) {
		// Return a copy so callers cannot mutate cached data.
		return append([]model.NodeData{}, f.cache...), true
	}

	slog.Debug("cache MISS", "endpoints", len(f.endpoints))
	results := make([]model.NodeData, len(f.endpoints))
	var wg sync.WaitGroup
	for i, ep := range f.endpoints {
		wg.Add(1)
		go func(i int, ep config.Endpoint) {
			defer wg.Done()
			results[i] = f.fetchOne(ctx, ep)
		}(i, ep)
	}
	wg.Wait()

	f.cache = results
	f.expiresAt = time.Now().Add(f.cacheTTL)
	// Return a copy so callers cannot mutate cached data.
	return append([]model.NodeData{}, results...), false
}

// fetchRaw fetches and parses Prometheus text-format metrics from a single URL.
func (f *Fetcher) fetchRaw(ctx context.Context, label, url string) ([]parser.Sample, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
		zfsErr          error
	)

	if ep.SmartctlURL != "" {
		// Fetch ZFS and smartctl concurrently; smartctl failure is non-fatal.
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			zfsSamples, zfsErr = f.fetchRaw(ctx, ep.Label, ep.URL)
		}()
		go func() {
			defer wg.Done()
			var err error
			smartctlSamples, err = f.fetchRaw(ctx, ep.Label, ep.SmartctlURL)
			if err != nil {
				slog.Warn("smartctl fetch failed (disk data unavailable)", "label", ep.Label, "error", err)
			}
		}()
		wg.Wait()
	} else {
		zfsSamples, zfsErr = f.fetchRaw(ctx, ep.Label, ep.URL)
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

	allSamples := append(zfsSamples, smartctlSamples...)
	nd.Pools = model.ExtractPools(allSamples)
	nd.ExporterInfo = model.ExtractExporterInfo(allSamples)
	nd.Disks = model.ExtractDisks(allSamples)
	if len(smartctlSamples) > 0 {
		nd.SmartctlInfo = model.ExtractSmartctlInfo(smartctlSamples)
	}
	slog.Debug("extracted pools", "label", ep.Label, "count", len(nd.Pools), "disks", len(nd.Disks))
	return nd
}
