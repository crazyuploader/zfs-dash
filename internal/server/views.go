package server

import (
	"time"

	"github.com/crazyuploader/zfs-dash/internal/config"
	"github.com/crazyuploader/zfs-dash/internal/model"
)

// pageData carries the fields every page template needs (topbar/nav state).
// Page-specific data structs embed it.
type pageData struct {
	ActiveTab      string // "storage" | "system" | "history"
	HistoryEnabled bool
	StorageEnabled bool
	RefreshSecs    int
}

// newPageData builds the shared page fields for the given active tab.
func newPageData(
	tab string,
	cfg *config.Config,
	historyEnabled bool,
	nodes []model.NodeData,
) pageData {
	data := pageData{
		ActiveTab:      tab,
		HistoryEnabled: historyEnabled,
		RefreshSecs:    int(cfg.Refresh.Seconds()),
	}
	for _, node := range nodes {
		if node.Exporters.StorageVisible() {
			data.StorageEnabled = true
			break
		}
	}
	return data
}

// nodeView is the browser-facing subset of NodeData, used both for the
// page's inline JS and the /api/metrics response.
// Scrape URLs stay in config and are never exposed to browsers or API consumers.
type nodeView struct {
	Label        string                 `json:"label"`
	Location     string                 `json:"location,omitempty"`
	FetchedAt    time.Time              `json:"fetched_at"`
	Error        string                 `json:"error,omitempty"`
	Exporters    model.ExporterStatuses `json:"exporters"`
	ExporterInfo model.ExporterInfo     `json:"exporter_info,omitempty"`
	SmartctlInfo model.SmartctlInfo     `json:"smartctl_info,omitempty"`
	Pools        []model.Pool           `json:"pools"`
	Disks        []model.DiskInfo       `json:"disks,omitempty"`
	System       *model.SystemInfo      `json:"system,omitempty"`
}

// systemView is the /api/system response row for one endpoint with a
// node_exporter available or explicitly required.
type systemView struct {
	Label     string                 `json:"label"`
	Location  string                 `json:"location,omitempty"`
	FetchedAt time.Time              `json:"fetched_at"`
	Error     string                 `json:"error,omitempty"`
	System    *model.SystemInfo      `json:"system,omitempty"`
	Exporters model.ExporterStatuses `json:"exporters"`
	PoolCount int                    `json:"pool_count"`
	DiskCount int                    `json:"disk_count"`
}

// systemPageData is the data passed to the system page template.
type systemPageData struct {
	pageData
	Nodes     []systemView
	FetchedAt string

	// Fleet KPI aggregates
	TotalNodes     int
	ExporterErrors int
	TotalCores     int
	AvgCPUPct      float64
	HasCPU         bool
	MemUsedBytes   float64
	MemTotal       float64
	MaxTempC       float64
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
		if v.Exporters.HasErrors() {
			d.ExporterErrors++
		}
		if v.System == nil {
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

// systemViews includes available and required node exporters in /api/system.
func systemViews(nodes []model.NodeData) []systemView {
	out := make([]systemView, 0, len(nodes))
	for _, node := range nodes {
		if node.Exporters.Node.Visible() {
			out = append(out, hostView(node))
		}
	}
	return out
}

// hostViews keeps all configured hosts on the homepage, including hosts for
// which automatic discovery has not found any exporters yet.
func hostViews(nodes []model.NodeData) []systemView {
	out := make([]systemView, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, hostView(node))
	}
	return out
}

func hostView(node model.NodeData) systemView {
	return systemView{
		Label:     node.Label,
		Location:  node.Location,
		FetchedAt: node.FetchedAt,
		Error:     node.Exporters.Node.Error,
		System:    node.System,
		Exporters: node.Exporters,
		PoolCount: len(node.Pools),
		DiskCount: len(node.Disks),
	}
}

// nodeViews converts fetched node data into its URL-stripped view form.
func nodeViews(nodes []model.NodeData) []nodeView {
	views := make([]nodeView, len(nodes))
	for i, n := range nodes {
		views[i] = nodeView{
			Label:        n.Label,
			Location:     n.Location,
			FetchedAt:    n.FetchedAt,
			Error:        n.Error,
			Exporters:    n.Exporters,
			ExporterInfo: n.ExporterInfo,
			SmartctlInfo: n.SmartctlInfo,
			Pools:        n.Pools,
			Disks:        n.Disks,
			System:       n.System,
		}
	}
	return views
}
