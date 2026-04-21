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
	endpoints []config.Endpoint
}

// New creates a Fetcher for the provided endpoints.
func New(endpoints []config.Endpoint) *Fetcher {
	return &Fetcher{
		client:    &http.Client{Timeout: fetchTimeout},
		endpoints: endpoints,
	}
}

// FetchAll fetches all endpoints concurrently and returns results in the same order.
func (f *Fetcher) FetchAll(ctx context.Context) []model.NodeData {
	results := make([]model.NodeData, len(f.endpoints))
	var wg sync.WaitGroup
	for i, ep := range f.endpoints {
		wg.Go(func() {
			results[i] = f.fetchOne(ctx, ep)
		})
	}
	wg.Wait()
	return results
}

func (f *Fetcher) fetchOne(ctx context.Context, ep config.Endpoint) model.NodeData {
	nd := model.NodeData{
		Label:     ep.Label,
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
		return nd
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		nd.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
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
	samples, err := parser.Parse(bytes.NewReader(body))
	if err != nil {
		nd.Error = fmt.Sprintf("parse: %v", err)
		return nd
	}
	nd.Pools = model.ExtractPools(samples)
	return nd
}
