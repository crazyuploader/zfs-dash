# zfs-dash

Minimal, real-time ZFS pool dashboard. Pulls metrics from one or more
[zfs_exporter](https://github.com/pdf/zfs_exporter) Prometheus endpoints and
renders a beautiful dashboard in your browser.

## Quick Start

```bash
# Install dependencies
go mod tidy

# Run with a config file
cp zfs-dash.yaml.example zfs-dash.yaml  # edit your endpoints
go run . serve

# Or pass endpoints directly
go run . serve \
  --endpoints http://proxmox-1:9134/metrics \
  --endpoints http://nas-1:9134/metrics \
  --addr :8080 \
  --refresh 15

# Build a static binary
go build -o zfs-dash .
./zfs-dash serve
```

Open http://localhost:8080

## Configuration

Config file: `./zfs-dash.yaml` (or `~/.config/zfs-dash/zfs-dash.yaml`)

```yaml
addr: ":8080"
refresh: 30

endpoints:
  - url: "http://host:9134/metrics"
    label: "my-server"
```

All config keys are also available as:

- Flags: `--addr`, `--refresh`, `--endpoints`
- Env vars: `ZFSDASH_ADDR`, `ZFSDASH_REFRESH`, `ZFSDASH_ENDPOINTS`

## API

| Route              | Description                   |
| ------------------ | ----------------------------- |
| `GET /`            | Dashboard HTML (SSR)          |
| `GET /api/metrics` | Raw JSON of all fetched pools |

## Directory Structure

```
zfs-dash/
├── main.go
├── go.mod
├── zfs-dash.yaml          ← config
├── cmd/
│   ├── root.go            ← cobra root + viper binding
│   └── serve.go           ← serve sub-command
├── internal/
│   ├── config/config.go   ← Config struct + Load()
│   ├── parser/parser.go   ← Prometheus text-format parser
│   ├── fetcher/fetcher.go ← concurrent HTTP fetcher
│   ├── model/model.go     ← domain types + ExtractPools()
│   └── server/server.go   ← Fiber v3 routes + template rendering
└── templates/
    ├── templates.go        ← re-exports dashboardHTML
    └── dashboard_html.go   ← full HTML template string
```
