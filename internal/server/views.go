package server

import (
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/crazyuploader/zfs-dash/internal/config"
	"github.com/crazyuploader/zfs-dash/internal/model"
)

// pageData carries the fields every page template needs (topbar/nav state).
// Page-specific data structs embed it.
type pageData struct {
	ActiveTab      string // "pools" | "system" | "history"
	HistoryEnabled bool
	SystemEnabled  bool
	RefreshSecs    int
}

// newPageData builds the shared page fields for the given active tab.
func newPageData(tab string, cfg *config.Config, historyEnabled bool) pageData {
	return pageData{
		ActiveTab:      tab,
		HistoryEnabled: historyEnabled,
		SystemEnabled:  systemConfigured(cfg),
		RefreshSecs:    int(cfg.Refresh.Seconds()),
	}
}

// systemConfigured reports whether any endpoint has a node_exporter URL.
func systemConfigured(cfg *config.Config) bool {
	return slices.ContainsFunc(cfg.Endpoints, func(ep config.Endpoint) bool {
		return ep.NodeExporterURL != ""
	})
}

// sanitizeError removes the scrape URL (and its host:port, which net errors
// embed separately, e.g. "dial tcp host:port") from fetch error messages so
// internal addresses are not exposed to browsers or API consumers.
func sanitizeError(msg, rawURL string) string {
	if msg == "" || rawURL == "" {
		return msg
	}
	msg = strings.ReplaceAll(msg, `"`+rawURL+`"`, "endpoint")
	msg = strings.ReplaceAll(msg, rawURL, "endpoint")
	if u, err := url.Parse(rawURL); err == nil && u.Host != "" {
		msg = strings.ReplaceAll(msg, u.Host, "endpoint")
	}
	return msg
}

// nodeView is the browser-facing subset of NodeData, used both for the
// page's inline JS and the /api/metrics response.
// URL is intentionally excluded so internal scrape endpoints are never
// exposed to browsers or API consumers.
type nodeView struct {
	Label        string             `json:"label"`
	Location     string             `json:"location,omitempty"`
	FetchedAt    time.Time          `json:"fetched_at"`
	Error        string             `json:"error,omitempty"`
	ExporterInfo model.ExporterInfo `json:"exporter_info,omitempty"`
	SmartctlInfo model.SmartctlInfo `json:"smartctl_info,omitempty"`
	Pools        []model.Pool       `json:"pools"`
	Disks        []model.DiskInfo   `json:"disks,omitempty"`
	System       *model.SystemInfo  `json:"system,omitempty"`
}

// systemView is the /api/system response row for one endpoint with a
// node_exporter configured.
type systemView struct {
	Label     string            `json:"label"`
	Location  string            `json:"location,omitempty"`
	FetchedAt time.Time         `json:"fetched_at"`
	Error     string            `json:"error,omitempty"`
	System    *model.SystemInfo `json:"system,omitempty"`
}

// systemPageData is the data passed to the system page template.
type systemPageData struct {
	pageData
	Nodes     []systemView
	FetchedAt string

	// Fleet KPI aggregates
	TotalNodes   int
	Unreachable  int
	TotalCores   int
	AvgCPUPct    float64
	HasCPU       bool
	MemUsedBytes float64
	MemTotal     float64
	MaxTempC     float64
}

// buildSystemPageData aggregates fleet KPIs over the system views.
func buildSystemPageData(views []systemView) systemPageData {
	d := systemPageData{
		Nodes:      views,
		FetchedAt:  time.Now().Format("15:04:05"),
		TotalNodes: len(views),
	}
	var cpuSum float64
	var cpuN int
	for _, v := range views {
		if v.System == nil {
			d.Unreachable++
			continue
		}
		s := v.System
		d.TotalCores += s.Cores
		d.MemTotal += s.MemTotal
		d.MemUsedBytes += s.MemTotal - s.MemAvailable
		if s.HasCPURates {
			cpuSum += s.CPUBusyPct
			cpuN++
		}
		for _, t := range s.Temps {
			if t.Celsius > d.MaxTempC {
				d.MaxTempC = t.Celsius
			}
		}
	}
	if cpuN > 0 {
		d.AvgCPUPct = cpuSum / float64(cpuN)
		d.HasCPU = true
	}
	return d
}

// systemViews returns view rows for nodes that have system data (or where a
// node_exporter is configured but returned nothing, so the UI can show an
// error card).
func systemViews(nodes []model.NodeData, cfg *config.Config) []systemView {
	configured := make(map[string]bool, len(cfg.Endpoints))
	for _, ep := range cfg.Endpoints {
		if ep.NodeExporterURL != "" {
			configured[ep.Label] = true
		}
	}
	// make (not var) so an empty result marshals as [] rather than null.
	out := make([]systemView, 0, len(nodes))
	for _, n := range nodes {
		if n.System == nil && !configured[n.Label] {
			continue
		}
		errMsg := ""
		if n.System == nil {
			errMsg = "node_exporter unreachable or returned no data"
		}
		out = append(out, systemView{
			Label:     n.Label,
			Location:  n.Location,
			FetchedAt: n.FetchedAt,
			Error:     errMsg,
			System:    n.System,
		})
	}
	return out
}

// nodeViews converts fetched node data into its URL-stripped view form.
func nodeViews(nodes []model.NodeData) []nodeView {
	views := make([]nodeView, len(nodes))
	for i, n := range nodes {
		views[i] = nodeView{
			Label:        n.Label,
			Location:     n.Location,
			FetchedAt:    n.FetchedAt,
			Error:        sanitizeError(n.Error, n.URL),
			ExporterInfo: n.ExporterInfo,
			SmartctlInfo: n.SmartctlInfo,
			Pools:        n.Pools,
			Disks:        n.Disks,
			System:       n.System,
		}
	}
	return views
}
