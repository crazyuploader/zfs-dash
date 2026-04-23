// Package fetcher fetches Prometheus metrics from multiple endpoints concurrently.
package fetcher

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
	Debug     bool
	cacheTTL  time.Duration
	cache     []model.NodeData
	expiresAt time.Time
}

// New creates a Fetcher for the provided endpoints.
func New(endpoints []config.Endpoint, debug bool, cacheTTL time.Duration) *Fetcher {
	return &Fetcher{
		client:    &http.Client{Timeout: fetchTimeout},
		endpoints: endpoints,
		Debug:     debug,
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
		if f.Debug {
			fmt.Printf("DEBUG: cache HIT (expires in %v)\n", time.Until(f.expiresAt).Round(time.Second))
		}
		data := f.cache
		f.mu.RUnlock()
		return data, true
	}
	f.mu.RUnlock()

	// Single-flight like behavior could be added here, but for now just standard lock.
	f.mu.Lock()
	defer f.mu.Unlock()

	// Re-check after acquiring write lock
	if time.Now().Before(f.expiresAt) {
		return f.cache, true
	}

	if f.Debug {
		fmt.Printf("DEBUG: cache MISS, fetching metrics from %d endpoints\n", len(f.endpoints))
	}
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
	return results, false
}

func (f *Fetcher) fetchOne(ctx context.Context, ep config.Endpoint) model.NodeData {
	if f.Debug {
		fmt.Printf("DEBUG: fetching from %s (%s)\n", ep.Label, ep.URL)
	}
	nd := model.NodeData{
		Label:     ep.Label,
		Location:  ep.Location,
		URL:       ep.URL,
		FetchedAt: time.Now(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep.URL, nil)
	if err != nil {
		nd.Error = fmt.Sprintf("build request: %v", err)
		return nd
	}
	resp, err := f.client.Do(req)
	if err != nil {
		nd.Error = fmt.Sprintf("unreachable: %v", err)
		if f.Debug {
			fmt.Printf("DEBUG: %s unreachable: %v\n", ep.Label, err)
		}
		return nd
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		nd.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
		if f.Debug {
			fmt.Printf("DEBUG: %s returned HTTP %d\n", ep.Label, resp.StatusCode)
		}
		return nd
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil {
		nd.Error = fmt.Sprintf("read: %v", err)
		return nd
	}
	if len(body) > maxResponseBytes {
		nd.Error = fmt.Sprintf("response too large: limit %d bytes", maxResponseBytes)
		return nd
	}
	if f.Debug {
		fmt.Printf("DEBUG: %s read %d bytes\n", ep.Label, len(body))
	}
	samples, err := parser.Parse(bytes.NewReader(body))
	if err != nil {
		nd.Error = fmt.Sprintf("parse: %v", err)
		if f.Debug {
			fmt.Printf("DEBUG: %s parse error: %v\n", ep.Label, err)
		}
		return nd
	}
	nd.Pools = model.ExtractPools(samples)
	if f.Debug {
		fmt.Printf("DEBUG: %s extracted %d pools\n", ep.Label, len(nd.Pools))
	}
	return nd
}
