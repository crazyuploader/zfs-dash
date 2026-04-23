package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/crazyuploader/zfs-dash/internal/config"
)

func TestFetcher_Cache(t *testing.T) {
	var callCount int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&callCount, 1)
		_, _ = fmt.Fprintln(w, "zfs_pool_health{pool=\"tank\"} 0")
	}))
	defer server.Close()

	eps := []config.Endpoint{
		{URL: server.URL, Label: "test-node"},
	}

	// Cache TTL of 1 second
	f := New(eps, false, 1*time.Second)
	ctx := context.Background()

	// First call - should hit the server
	_, isCached := f.FetchAll(ctx)
	if isCached {
		t.Fatal("expected first call to be a cache miss")
	}
	if atomic.LoadInt64(&callCount) != 1 {
		t.Fatalf("expected 1 call, got %d", atomic.LoadInt64(&callCount))
	}

	// Second call - should be cached
	_, isCached = f.FetchAll(ctx)
	if !isCached {
		t.Fatal("expected second call to be a cache hit")
	}
	if atomic.LoadInt64(&callCount) != 1 {
		t.Fatalf("expected cached result, but server was hit again (call count: %d)", atomic.LoadInt64(&callCount))
	}

	// Wait for cache to expire
	time.Sleep(1100 * time.Millisecond)

	// Third call - should hit the server again
	_, isCached = f.FetchAll(ctx)
	if isCached {
		t.Fatal("expected call after expiry to be a cache miss")
	}
	if atomic.LoadInt64(&callCount) != 2 {
		t.Fatalf("expected 2 calls after expiry, got %d", atomic.LoadInt64(&callCount))
	}
}

func TestFetcher_SetEndpoints(t *testing.T) {
	var callCount int64
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&callCount, 1)
		_, _ = fmt.Fprintln(w, "zfs_pool_health{pool=\"tank1\"} 0")
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&callCount, 1)
		_, _ = fmt.Fprintln(w, "zfs_pool_health{pool=\"tank2\"} 0")
	}))
	defer server2.Close()

	f := New([]config.Endpoint{{URL: server1.URL, Label: "node1"}}, false, 1*time.Minute)
	ctx := context.Background()

	// Initial fetch
	nodes, isCached := f.FetchAll(ctx)
	if isCached {
		t.Fatal("expected first fetch to be a cache miss")
	}
	if len(nodes) != 1 || nodes[0].Label != "node1" {
		t.Fatalf("expected node1, got %+v", nodes)
	}

	// Update endpoints
	f.SetEndpoints([]config.Endpoint{{URL: server2.URL, Label: "node2"}})

	// Fetch again - should hit server2 and ignore old cache
	nodes, isCached = f.FetchAll(ctx)
	if isCached {
		t.Fatal("expected fetch after SetEndpoints to be a cache miss")
	}
	if len(nodes) != 1 || nodes[0].Label != "node2" {
		t.Fatalf("expected node2 after hot-reload, got %+v", nodes)
	}
}
