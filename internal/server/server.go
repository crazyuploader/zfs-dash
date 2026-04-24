// Package server boots the Fiber v3 web server and serves the ZFS dashboard.
package server

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/crazyuploader/zfs-dash/internal/config"
	"github.com/crazyuploader/zfs-dash/internal/fetcher"
	"github.com/crazyuploader/zfs-dash/internal/history"
	"github.com/crazyuploader/zfs-dash/internal/model"
	"github.com/crazyuploader/zfs-dash/templates"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

// Hub broadcasts reload signals to connected SSE clients.
// clients is a sync.Map keyed by chan bool.
type Hub struct {
	clients sync.Map
}

func newHub() *Hub {
	return &Hub{}
}

func (h *Hub) broadcast() {
	h.clients.Range(func(key, value any) bool {
		ch := key.(chan bool)
		select {
		case ch <- true:
		default:
			h.clients.Delete(ch)
		}
		return true
	})
}

const (
	httpReadTimeout    = 15 * time.Second
	httpIdleTimeout    = 60 * time.Second
	httpHandlerTimeout = 15 * time.Second
)

// nodeView is the subset of NodeData serialized into the page's inline JS.
// URL is intentionally excluded to avoid exposing internal scrape endpoints to browsers.
type nodeView struct {
	Label        string             `json:"label"`
	Location     string             `json:"location,omitempty"`
	ExporterInfo model.ExporterInfo `json:"exporter_info,omitempty"`
	SmartctlInfo model.SmartctlInfo `json:"smartctl_info,omitempty"`
	Pools        []model.Pool       `json:"pools"`
	Disks        []model.DiskInfo   `json:"disks,omitempty"`
}

// templateData is the data passed to the HTML template.
type templateData struct {
	Nodes            []model.NodeData
	NodesJSON        template.JS // URL-stripped JSON for inline script
	RefreshSecs      int
	FetchedAt        string
	TotalPools       int
	UnreachableNodes int
	HealthyPools     int
	DegradedPools    int
	ErroredPools     int
	TotalNodes       int
	HistoryEnabled   bool
}

// Start registers routes and begins listening.
func Start(cfg *config.Config) error {
	setupLogger(cfg)

	var cfgPtr atomic.Pointer[config.Config]
	cfgPtr.Store(cfg)

	slog.Debug("starting server in debug mode", "config", cfg)

	f := fetcher.New(cfg.Endpoints, cfg.CacheTTL)
	hub := newHub()

	// Graceful shutdown context — cancelled on SIGTERM/SIGINT.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdownSigs := make(chan os.Signal, 1)
	signal.Notify(shutdownSigs, syscall.SIGTERM, os.Interrupt)

	// Hot-reload config on SIGHUP
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)
	go func() {
		for range sigs {
			slog.Info("SIGHUP received, reloading config...")
			newCfg, err := config.Load()
			if err != nil {
				slog.Error("config reload failed", "error", err)
				continue
			}
			setupLogger(newCfg)
			f.SetEndpoints(newCfg.Endpoints)
			cfgPtr.Store(newCfg)
			slog.Info("config reloaded successfully")
		}
	}()

	// History store and recorder (optional).
	var histStore *history.Store
	if cfg.History.Enabled {
		var err error
		histStore, err = history.Open(cfg.History.Path, cfg.History.Retention)
		if err != nil {
			slog.Error("history store failed to open — history disabled", "error", err, "path", cfg.History.Path)
			cfg.History.Enabled = false
		} else {
			defer func() { _ = histStore.Close() }()
			recInterval := cfg.History.RecordInterval
			if recInterval <= 0 {
				recInterval = cfg.Refresh
			}
			rec := history.NewRecorder(histStore, f, recInterval)
			rec.Start(ctx)
			slog.Info("history enabled", "path", cfg.History.Path, "retention", cfg.History.Retention, "record_interval", recInterval)
		}
	}

	tmpl, err := template.New("dashboard").Funcs(funcMap()).Parse(templates.Dashboard)
	if err != nil {
		return fmt.Errorf("template parse: %w", err)
	}

	var histTmpl *template.Template
	if histStore != nil {
		histTmpl, err = template.New("history").Funcs(funcMap()).Parse(templates.History)
		if err != nil {
			return fmt.Errorf("history template parse: %w", err)
		}
	}

	app := fiber.New(fiber.Config{
		AppName:      "zfs-dash",
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: 0, // Disable write timeout for SSE streams
		IdleTimeout:  httpIdleTimeout,
		TrustProxy:   len(cfg.TrustedProxies) > 0,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Proxies: cfg.TrustedProxies,
		},
		ProxyHeader: fiber.HeaderXForwardedFor,
	})

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${ips} ${method} ${path}\n",
	}))

	app.Use(func(c fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Content-Security-Policy", "default-src 'self'; script-src 'unsafe-inline'; style-src 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'")
		return c.Next()
	})

	rl := limiter.New(limiter.Config{
		Max:        60,
		Expiration: 1 * time.Minute,
	})

	// SSE Endpoint
	app.Get("/events", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		notify := make(chan bool, 1)
		hub.clients.Store(notify, true)
		defer hub.clients.Delete(notify)

		clientIP := c.IP()
		c.Response().SetBodyStreamWriter(func(w *bufio.Writer) {
			slog.Debug("SSE client connected", "ip", clientIP)

			// Send initial keep-alive
			_, _ = fmt.Fprintf(w, ":\n\n")
			_ = w.Flush()

			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-notify:
					_, _ = fmt.Fprintf(w, "data: refresh\n\n")
					if err := w.Flush(); err != nil {
						return
					}
				case <-c.Context().Done():
					slog.Debug("SSE client disconnected", "ip", clientIP)
					return
				case <-ticker.C:
					// keep-alive
					_, _ = fmt.Fprintf(w, ":\n\n")
					if err := w.Flush(); err != nil {
						return
					}
				}
			}
		})

		return nil
	})

	// JSON API
	app.Get("/api/metrics", rl, func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), httpHandlerTimeout)
		defer cancel()
		nodes, isCached := f.FetchAll(ctx)
		if !isCached {
			hub.broadcast()
		}
		setCacheHeaders(c, f, isCached)
		return c.JSON(nodes)
	})

	app.Get("/api/health/:label", rl, func(c fiber.Ctx) error {
		curCfg := cfgPtr.Load()
		return serveHealthCheck(c, f, c.Params("label"), "", curCfg)
	})

	app.Get("/api/health/:label/:pool", rl, func(c fiber.Ctx) error {
		curCfg := cfgPtr.Load()
		return serveHealthCheck(c, f, c.Params("label"), c.Params("pool"), curCfg)
	})

	// Dashboard
	app.Get("/", func(c fiber.Ctx) error {
		curCfg := cfgPtr.Load()
		reqCtx, cancel := context.WithTimeout(c.Context(), httpHandlerTimeout)
		defer cancel()

		nodes, isCached := f.FetchAll(reqCtx)
		if !isCached {
			hub.broadcast()
		}
		data := buildTemplateData(nodes, curCfg)
		data.HistoryEnabled = histStore != nil

		setCacheHeaders(c, f, isCached)
		c.Set("Cache-Control", "no-store")
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			slog.Error("template execution failed", "error", err)
			return fiber.ErrInternalServerError
		}

		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.Send(buf.Bytes())
	})

	// History page + API (only when history is enabled)
	if histStore != nil {
		app.Get("/history", func(c fiber.Ctx) error {
			curCfg := cfgPtr.Load()
			var buf bytes.Buffer
			if err := histTmpl.Execute(&buf, map[string]any{"RefreshSecs": int(curCfg.Refresh.Seconds())}); err != nil {
				slog.Error("history template execution failed", "error", err)
				return fiber.ErrInternalServerError
			}
			c.Set("Content-Type", "text/html; charset=utf-8")
			c.Set("Cache-Control", "no-store")
			return c.Send(buf.Bytes())
		})

		app.Get("/api/history/series", rl, func(c fiber.Ctx) error {
			series, err := histStore.ListSeries()
			if err != nil {
				slog.Error("history list series failed", "error", err)
				return fiber.ErrInternalServerError
			}
			if series == nil {
				series = []history.SeriesInfo{}
			}
			return c.JSON(series)
		})

		app.Get("/api/history/query", rl, func(c fiber.Ctx) error {
			key := c.Query("key")
			if key == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "key required"})
			}
			var fromUnix, toUnix, bucketSecs int64
			if s := c.Query("from"); s != "" {
				v, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid from"})
				}
				fromUnix = v
			}
			if s := c.Query("to"); s != "" {
				v, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid to"})
				}
				toUnix = v
			}
			if s := c.Query("bucket", "0"); s != "" {
				v, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid bucket"})
				}
				bucketSecs = v
			}

			if fromUnix < 0 || toUnix < 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "timestamps must be non-negative"})
			}
			now := time.Now()
			to := now
			if toUnix > 0 {
				to = time.Unix(toUnix, 0)
			}
			from := to.Add(-24 * time.Hour) // default: 24h before to
			if fromUnix > 0 {
				from = time.Unix(fromUnix, 0)
			}

			points, err := histStore.Query(key, from, to, bucketSecs)
			if err != nil {
				slog.Error("history query failed", "error", err, "key", key)
				return fiber.ErrInternalServerError
			}
			if points == nil {
				points = []history.Point{}
			}
			return c.JSON(points)
		})
	}

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Shutdown on SIGTERM/SIGINT
	go func() {
		<-shutdownSigs
		slog.Info("shutdown signal received")
		cancel()
		_ = app.Shutdown()
	}()

	slog.Info("zfs-dash started", "url", fmt.Sprintf("http://localhost%s", cfg.Addr))
	return app.Listen(cfg.Addr)
}

func setupLogger(cfg *config.Config) {
	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}

	var handler slog.Handler
	if cfg.LogFormat == "json" {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	}

	slog.SetDefault(slog.New(handler))
}

func buildTemplateData(nodes []model.NodeData, cfg *config.Config) templateData {
	views := make([]nodeView, len(nodes))
	for i, n := range nodes {
		views[i] = nodeView{Label: n.Label, Location: n.Location, ExporterInfo: n.ExporterInfo, SmartctlInfo: n.SmartctlInfo, Pools: n.Pools, Disks: n.Disks}
	}
	nodesJSON, _ := json.Marshal(views)

	d := templateData{
		Nodes:       nodes,
		NodesJSON:   template.JS(nodesJSON), //nolint:gosec // json.Marshal output is safe for inline JS
		RefreshSecs: int(cfg.Refresh.Seconds()),
		FetchedAt:   time.Now().Format("15:04:05"),
		TotalNodes:  len(nodes),
		// HistoryEnabled is set by the caller after buildTemplateData returns.
	}
	for _, n := range nodes {
		if n.Error != "" {
			d.UnreachableNodes++
		}
		for _, p := range n.Pools {
			d.TotalPools++
			switch p.Health {
			case model.HealthOnline:
				d.HealthyPools++
			case model.HealthDegraded:
				d.DegradedPools++
			default:
				d.ErroredPools++
			}
		}
	}
	return d
}

func serveHealthCheck(c fiber.Ctx, f *fetcher.Fetcher, label, poolName string, cfg *config.Config) error {
	slog.Debug("health check", "label", label, "pool", poolName)

	ctx, cancel := context.WithTimeout(c.Context(), httpHandlerTimeout)
	defer cancel()

	nodes, isCached := f.FetchAll(ctx)
	setCacheHeaders(c, f, isCached)

	node, err := findNodeByLabel(nodes, label)
	if err != nil {
		slog.Debug("node not found", "label", label)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "not_found",
			"label":  label,
		})
	}

	if node.Error != "" {
		slog.Debug("node has error", "label", label, "error", node.Error)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":   "down",
			"label":    node.Label,
			"location": node.Location,
		})
	}

	if poolName == "" {
		badPools := []string{}
		overThreshold := []string{}
		for _, pool := range node.Pools {
			if pool.Health != model.HealthOnline {
				badPools = append(badPools, pool.Name)
			} else if cfg.MaxUsagePercent > 0 && pool.UsedPercent > cfg.MaxUsagePercent {
				overThreshold = append(overThreshold, pool.Name)
			}
		}

		status := fiber.StatusOK
		state := "up"
		reason := ""
		if len(node.Pools) == 0 {
			status = fiber.StatusServiceUnavailable
			state = "no_pools"
			slog.Debug("node has 0 pools", "label", label)
		} else if len(badPools) > 0 {
			status = fiber.StatusServiceUnavailable
			state = "degraded"
			reason = "unhealthy_pools"
			slog.Debug("node has unhealthy pools", "label", label, "pools", badPools)
		} else if len(overThreshold) > 0 {
			status = fiber.StatusServiceUnavailable
			state = "degraded"
			reason = "pool_over_threshold"
			slog.Debug("node has pools over threshold", "label", label, "pools", overThreshold, "threshold", cfg.MaxUsagePercent)
		}

		res := fiber.Map{
			"status":          state,
			"label":           node.Label,
			"location":        node.Location,
			"pool_count":      len(node.Pools),
			"unhealthy_pools": badPools,
		}
		if reason != "" {
			res["reason"] = reason
		}
		if len(overThreshold) > 0 {
			res["over_threshold_pools"] = overThreshold
		}

		return c.Status(status).JSON(res)
	}

	pool, err := findPoolByName(node.Pools, poolName)
	if err != nil {
		slog.Debug("pool not found", "label", label, "pool", poolName)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "down",
			"label":  node.Label,
			"pool":   poolName,
			"reason": "pool_not_found",
		})
	}

	status := fiber.StatusOK
	state := "up"
	reason := ""
	if pool.Health != model.HealthOnline {
		status = fiber.StatusServiceUnavailable
		state = "degraded"
		reason = "pool_unhealthy"
		slog.Debug("pool health is not ONLINE", "label", label, "pool", poolName, "health", pool.Health)
	} else if cfg.MaxUsagePercent > 0 && pool.UsedPercent > cfg.MaxUsagePercent {
		status = fiber.StatusServiceUnavailable
		state = "degraded"
		reason = "pool_over_threshold"
		slog.Debug("pool is over threshold", "label", label, "pool", poolName, "used_percent", pool.UsedPercent, "threshold", cfg.MaxUsagePercent)
	}

	res := fiber.Map{
		"status":   state,
		"label":    node.Label,
		"location": node.Location,
		"pool":     pool.Name,
		"health":   pool.Health,
	}
	if reason != "" {
		res["reason"] = reason
	}
	if status != fiber.StatusOK {
		res["used_percent"] = pool.UsedPercent
	}

	return c.Status(status).JSON(res)
}

func findNodeByLabel(nodes []model.NodeData, label string) (*model.NodeData, error) {
	for i := range nodes {
		if nodes[i].Label == label {
			return &nodes[i], nil
		}
	}
	return nil, fmt.Errorf("label %q not found", label)
}

func findPoolByName(pools []model.Pool, name string) (*model.Pool, error) {
	for i := range pools {
		if pools[i].Name == name {
			return &pools[i], nil
		}
	}
	return nil, fmt.Errorf("pool %q not found", name)
}

func setCacheHeaders(c fiber.Ctx, f *fetcher.Fetcher, isCached bool) {
	if isCached {
		c.Set("X-Cache", "HIT")
		expiresAt, _ := f.CacheInfo()
		if time.Now().Before(expiresAt) {
			c.Set("X-Cache-Expires-In", time.Until(expiresAt).Round(time.Second).String())
		}
	} else {
		c.Set("X-Cache", "MISS")
	}
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"humanBytes": model.HumanBytes,
		"healthClass": func(h model.PoolHealth) string {
			switch h {
			case model.HealthOnline:
				return "health-online"
			case model.HealthDegraded:
				return "health-degraded"
			default:
				return "health-faulted"
			}
		},
		"fmtNodeTime": func(t time.Time) string {
			return t.Format("15:04:05")
		},
		"toJSON": func(v any) string {
			b, _ := json.Marshal(v)
			return string(b)
		},
		"fmtSpeed": func(bps float64) string {
			switch {
			case bps >= 1e9:
				return fmt.Sprintf("%.0f Gb/s", bps/1e9)
			case bps >= 1e6:
				return fmt.Sprintf("%.0f Mb/s", bps/1e6)
			default:
				return fmt.Sprintf("%.0f b/s", bps)
			}
		},
		"exitStatusDesc": func(code float64) string {
			n := int(code)
			if n == 0 {
				return ""
			}
			var parts []string
			if n&(1<<1) != 0 {
				parts = append(parts, "device failure")
			}
			if n&(1<<2) != 0 {
				parts = append(parts, "disk failing")
			}
			if n&(1<<3) != 0 {
				parts = append(parts, "prefail attributes")
			}
			if n&(1<<4) != 0 {
				parts = append(parts, "prev failed attributes")
			}
			if n&(1<<5) != 0 {
				parts = append(parts, "error log has errors")
			}
			if n&(1<<6) != 0 {
				parts = append(parts, "self-test errors")
			}
			if len(parts) == 0 {
				return fmt.Sprintf("code %d", n)
			}
			return strings.Join(parts, ", ")
		},
		"diskHasIssues": func(d model.DiskInfo) bool {
			return d.PendingSectors > 0 || d.OfflineUncorrectable > 0 || d.ReportedUncorrect > 0 ||
				d.ProgramFailCount > 0 || d.EraseFailCount > 0 ||
				(d.HasExitStatus && d.ExitStatus > 0)
		},
		"tempBarPct": func(temp, maxTemp float64) string {
			if maxTemp <= 0 {
				maxTemp = 70
			}
			pct := (temp / maxTemp) * 100
			if pct > 100 {
				pct = 100
			} else if pct < 0 {
				pct = 0
			}
			return fmt.Sprintf("%.1f", pct)
		},
		"fmtHours": func(h float64) string {
			total := int(h)
			days := total / 24
			hrs := total % 24
			if days >= 365 {
				y := days / 365
				d := days % 365
				return fmt.Sprintf("%dy %dd", y, d)
			}
			if days > 0 {
				return fmt.Sprintf("%dd %dh", days, hrs)
			}
			return fmt.Sprintf("%dh", total)
		},
		"maskSerial": func(s string) string {
			const maskLen = 5
			if len(s) <= maskLen {
				return strings.Repeat("x", len(s))
			}
			return s[:len(s)-maskLen] + strings.Repeat("x", maskLen)
		},
		"diskTypeLabel": func(iface string, rpm int) string {
			switch {
			case iface == "nvme":
				return "NVMe"
			case rpm > 0:
				return "HDD"
			default:
				return "SSD"
			}
		},
		"diskTypeClass": func(iface string, rpm int) string {
			switch {
			case iface == "nvme":
				return "nvme"
			case rpm > 0:
				return "hdd"
			default:
				return "ssd"
			}
		},
		"tempClass": func(c float64) string {
			switch {
			case c > 55:
				return "hot"
			case c > 45:
				return "warm"
			default:
				return "cool"
			}
		},
		"gt0":    func(f float64) bool { return f > 0 },
		"gte":    func(a, b float64) bool { return a >= b },
		"mul100": func(f float64) float64 { return f * 100 },
		"mul512": func(f float64) float64 { return f * 512 },
		"join":   strings.Join,
		"safeJS": func(s string) template.JS { return template.JS(s) },
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call: expected even number of arguments")
			}
			dict := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}
}
