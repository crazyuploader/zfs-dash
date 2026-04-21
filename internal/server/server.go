// Package server boots the Fiber v3 web server and serves the ZFS dashboard.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/yourname/zfs-dash/internal/config"
	"github.com/yourname/zfs-dash/internal/fetcher"
	"github.com/yourname/zfs-dash/internal/model"
	"github.com/yourname/zfs-dash/templates"
)

// templateData is the data passed to the HTML template.
// All fields are pre-computed so the template stays logic-free.
type templateData struct {
	Nodes         []model.NodeData
	RefreshSecs   int
	FetchedAt     string // human-readable timestamp of the current fetch
	TotalPools    int
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
		AppName: "zfs-dash",
	})

	// JSON API — useful for scripting / alerts.
	app.Get("/api/metrics", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
		defer cancel()
		nodes := f.FetchAll(ctx)
		return c.JSON(nodes)
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
		"gt0": func(f float64) bool { return f > 0 },
		"gte": func(a, b float64) bool { return a >= b },

		// join wraps strings.Join for use in templates.
		"join": strings.Join,

		// safeJS wraps a string as template.JS to skip escaping inside <script>.
		"safeJS": func(s string) template.JS { return template.JS(s) },
	}
}
