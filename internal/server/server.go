// Package server boots the Fiber v3 web server and serves the ZFS dashboard.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/crazyuploader/zfs-dash/internal/config"
	"github.com/crazyuploader/zfs-dash/internal/fetcher"
	"github.com/crazyuploader/zfs-dash/internal/model"
	"github.com/crazyuploader/zfs-dash/templates"
	"github.com/gofiber/fiber/v3"
)

const (
	httpReadTimeout  = 15 * time.Second
	httpWriteTimeout = 15 * time.Second
	httpIdleTimeout  = 60 * time.Second
)

// templateData is the data passed to the HTML template.
// All fields are pre-computed so the template stays logic-free.
type templateData struct {
	Nodes         []model.NodeData
	RefreshSecs   int
	FetchedAt     string // human-readable timestamp of the current fetch
	TotalPools    int
	UnreachableNodes int
	HealthyPools  int
	DegradedPools int
	ErroredPools  int
	TotalNodes    int
}

// Start registers routes and begins listening.
func Start(cfg *config.Config) error {
	f := fetcher.New(cfg.Endpoints)

	tmpl, err := template.New("dashboard").Funcs(funcMap()).Parse(templates.Dashboard)
	if err != nil {
		return fmt.Errorf("template parse: %w", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "zfs-dash",
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,
		IdleTimeout:  httpIdleTimeout,
	})

	// JSON API — useful for scripting / alerts.
	app.Get("/api/metrics", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
		defer cancel()
		nodes := f.FetchAll(ctx)
		return c.JSON(nodes)
	})

	app.Get("/api/health/:label", func(c fiber.Ctx) error {
		return serveHealthCheck(c, f, c.Params("label"), "")
	})

	app.Get("/api/health/:label/:pool", func(c fiber.Ctx) error {
		return serveHealthCheck(c, f, c.Params("label"), c.Params("pool"))
	})

	// Dashboard — SSR HTML page.
	app.Get("/", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
		defer cancel()

		nodes := f.FetchAll(ctx)
		data := buildTemplateData(nodes, cfg)

		c.Set("Content-Type", "text/html; charset=utf-8")
		return tmpl.Execute(c.Response().BodyWriter(), data)
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	fmt.Printf("→  zfs-dash  http://localhost%s\n", cfg.Addr)
	return app.Listen(cfg.Addr)
}

func buildTemplateData(nodes []model.NodeData, cfg *config.Config) templateData {
	d := templateData{
		Nodes:       nodes,
		RefreshSecs: int(cfg.Refresh.Seconds()),
		FetchedAt:   time.Now().Format("15:04:05"),
		TotalNodes:  len(nodes),
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

func serveHealthCheck(c fiber.Ctx, f *fetcher.Fetcher, label, poolName string) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	nodes := f.FetchAll(ctx)
	node, err := findNodeByLabel(nodes, label)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "not_found",
			"label":  label,
		})
	}

	if node.Error != "" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":   "down",
			"label":    node.Label,
			"location": node.Location,
		})
	}

	if poolName == "" {
		badPools := []string{}
		for _, pool := range node.Pools {
			if pool.Health != model.HealthOnline {
				badPools = append(badPools, pool.Name)
			}
		}

		status := fiber.StatusOK
		state := "up"
		if len(badPools) > 0 {
			status = fiber.StatusServiceUnavailable
			state = "degraded"
		}

		return c.Status(status).JSON(fiber.Map{
			"status":          state,
			"label":           node.Label,
			"location":        node.Location,
			"pool_count":      len(node.Pools),
			"unhealthy_pools": badPools,
		})
	}

	pool, err := findPoolByName(node.Pools, poolName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "not_found",
			"label":  node.Label,
			"pool":   poolName,
		})
		}

	status := fiber.StatusOK
	state := "up"
	if pool.Health != model.HealthOnline {
		status = fiber.StatusServiceUnavailable
		state = "degraded"
	}

	return c.Status(status).JSON(fiber.Map{
		"status":   state,
		"label":    node.Label,
		"location": node.Location,
		"pool":     pool.Name,
		"health":   pool.Health,
	})
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

func funcMap() template.FuncMap {
	return template.FuncMap{
		// humanBytes converts float64 bytes → "3.72 TB"
		"humanBytes": model.HumanBytes,

		// healthClass returns a CSS class string for the health badge.
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

		// fmtNodeTime formats the FetchedAt field on a NodeData.
		"fmtNodeTime": func(t time.Time) string {
			return t.Format("15:04:05")
		},

		// toJSON marshals any value to a JSON string (safe for use in <script>).
		"toJSON": func(v any) string {
			b, _ := json.Marshal(v)
			return string(b)
		},

		// gt2 is a two-arg greater-than helper (template's gt needs comparable types).
		"gt0":    func(f float64) bool { return f > 0 },
		"gte":    func(a, b float64) bool { return a >= b },
		"mul100": func(f float64) float64 { return f * 100 },

		// join wraps strings.Join for use in templates.
		"join": strings.Join,

		// safeJS wraps a string as template.JS to skip escaping inside <script>.
		"safeJS": func(s string) template.JS { return template.JS(s) },

		// dict creates a map from a list of alternating keys and values.
		// Useful for passing multiple arguments to a sub-template.
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
