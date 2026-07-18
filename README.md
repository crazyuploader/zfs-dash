# zfs-dash

Minimal dashboard for ZFS pool metrics from [pdf/zfs_exporter](https://github.com/pdf/zfs_exporter).
Two optional companions extend it: [smartctl_exporter](https://github.com/prometheus-community/smartctl_exporter) adds disk health, and [node_exporter](https://github.com/prometheus/node_exporter) adds a per-node System tab (CPU, memory, load, filesystems, network, temperatures).

## Run

```bash
cp config.yaml.example config.yaml
go run . serve --config config.yaml
```

Open `http://localhost:8054`.

### Flags

- `--config`: Config file (default: `./config.yaml`).
- `--addr`: Address to listen on (default: `:8054`).
- `--endpoints`: Comma-separated or repeated exporter URLs.
- `--refresh`: Auto-refresh interval in seconds (default: `300`). The config file also accepts duration strings like `"5m"`.
- `--debug`: Enable verbose debug logging.
- `--trusted-proxies`: List of trusted proxy IPs for reverse proxy header support.
- `--max-usage-percent`: Usage threshold for health failure (0 to disable).
- `--log-format`: Log format, either `text` or `json` (default: `text`).
- `--history-enabled`: Enable time-series history storage (default: `false`).
- `--history-path`: Path to history database file (default: `./data/history.db`).
- `--history-retention`: Retention period, e.g. `720h` for 30 days (default: `720h`).
- `--history-record-interval`: How often to record history samples, e.g. `5m` (default: same as `--refresh`).

## Config

```yaml
addr: ":8054"
refresh: 300
cache_ttl: 30 # cache fetched metrics for 30 seconds
max_usage_percent: 90 # fail health check if any pool > 90% full
log_format: "text" # "text" or "json"
debug: false
trusted_proxies: [] # e.g., ["127.0.0.1", "100.64.0.0/10"]

history:
  enabled: false
  path: "./data/history.db"
  retention: "720h" # 30 days; supports any Go duration string
  record_interval: "5m" # sample frequency; defaults to refresh interval

endpoints:
  - url: "http://host1:9134/metrics"
    label: "node-1"
    location: "Singapore"
    smartctl_url: "http://host1:9633/metrics" # optional: disk health
    node_exporter_url: "http://host1:9100/metrics" # optional: System tab
  - url: "http://host2:9134/metrics"
    label: "node-2"
    location: "Mumbai"
```

## System Tab

When any endpoint has a `node_exporter_url`, a **System** tab appears with per-node cards:

- CPU busy % and iowait (computed from consecutive scrapes; shows "warming up" until the second sample, ~30s after start)
- Pressure stall percentages (CPU / IO / memory PSI)
- Memory usage, buffers/cache, swap
- Load averages colored relative to core count
- Real (non-ZFS) filesystem usage; skips tmpfs, fuse, and overlay mounts
- Network throughput; skips loopback and Proxmox guest interfaces (`veth*`, `tap*`, `fwbr*`)
- hwmon temperature sensors
- OS, kernel, and uptime

JSON equivalent at `GET /api/system`.

## History

When `history.enabled: true`, zfs-dash records pool usage and disk metrics (temperature, wear, power-on hours) to a local [bbolt](https://github.com/etcd-io/bbolt) database, sampling every `history.record_interval` (defaults to `refresh`).

Charts live at **`/history`**; the History tab appears in the topbar once enabled.

**Recorded metrics:**

| Series                    | Description                  |
| ------------------------- | ---------------------------- |
| `pool/{name}/used_pct`    | Pool used %                  |
| `pool/{name}/alloc_bytes` | Pool allocated bytes         |
| `pool/{name}/free_bytes`  | Pool free bytes              |
| `disk/{dev}/temp_c`       | Disk temperature °C          |
| `disk/{dev}/wear_pct`     | NVMe percentage used (wear)  |
| `disk/{dev}/wear_lvl`     | SATA SSD wear leveling count |
| `disk/{dev}/pow_hrs`      | Power-on hours               |
| `system/node/cpu_pct`     | CPU busy % (node_exporter)   |
| `system/node/mem_used_pct`| Memory used % (node_exporter)|
| `system/node/load1`       | 1-minute load average        |

zfs-dash prunes data older than the retention window. Each data point stores 8 bytes of values: 30 days at a 5-minute interval across 50 disks × 4 metrics ≈ 14 MB raw, around 35 MB on disk with bbolt key and page overhead.

**Docker:** uncomment the `zfs-dash-data` volume in `docker-compose.yml` and set `history.path: /data/history.db` in your config.

## Hot Reload

Reload configuration (endpoints, debug mode, cache settings) without restarting:

```bash
kill -HUP $(pgrep zfs-dash)
```

## Docker

```bash
docker compose up -d
```

Edit `config.yaml` (copy from `config.yaml.example`) before starting the stack.

## API

- `GET /`
- `GET /system`: per-node system metrics (needs `node_exporter_url` on an endpoint)
- `GET /history`: history charts (requires `history.enabled: true`)
- `GET /api/metrics`
- `GET /api/system`: system metrics as JSON
- `GET /api/history/series`: list recorded series (history only)
- `GET /api/history/query?key=&from=&to=&bucket=`: query time-series data (history only)
- `GET /api/health/:label`
- `GET /api/health/:label/:pool`
- `GET /health`

The web UI and API responses never include scrape endpoint URLs. SSE connections on `/events` cap at 64 concurrent clients.

### Health Checks

Health checks suit monitoring tools like Uptime Kuma: `200 OK` when healthy, `503 Service Unavailable` otherwise.

- **Host Health:** `GET /api/health/:label`
  - Returns `503` if the host is unreachable or no ZFS pools are detected.
  - Returns `503` if any pool is degraded or faulted.
  - Returns `503` if any pool exceeds `max_usage_percent`.
- **Pool Health:** `GET /api/health/:label/:pool`
  - Returns `503` if the pool is not found or not healthy.
  - Returns `503` if the pool exceeds `max_usage_percent`.

Examples:

- `GET /api/health/node-1`
- `GET /api/health/node-1/tank`
