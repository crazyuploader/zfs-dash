package server

import (
	"net/url"
	"strings"
	"time"

	"github.com/crazyuploader/zfs-dash/internal/model"
)

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
		}
	}
	return views
}
