# zfs-dash

Minimal dashboard for ZFS pool metrics from [pdf/zfs_exporter](https://github.com/pdf/zfs_exporter).

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
- `--refresh`: Auto-refresh interval in seconds (default: `300`).
- `--debug`: Enable verbose debug logging.
- `--trusted-proxies`: List of trusted proxy IPs for reverse proxy header support.
- `--max-usage-percent`: Usage threshold for health failure (0 to disable).
- `--log-format`: Log format, either `text` or `json` (default: `text`).

## Config

```yaml
addr: ":8054"
refresh: 300
cache_ttl: 30 # Cache fetched metrics for 30 seconds to reduce load
max_usage_percent: 90 # Fail health check if any pool > 90% full
log_format: "text" # "text" or "json"
debug: false
trusted_proxies: [] # List of proxy IPs or CIDR ranges (e.g., ["127.0.0.1", "100.64.0.0/10"])
```
endpoints:
  - url: "http://host1:9134/metrics"
    label: "node-1"
    location: "Singapore"
```

## Hot Reload

You can reload the configuration (endpoints, debug mode, and cache settings) without restarting the server by sending a `SIGHUP` signal:

```bash
kill -HUP $(pgrep zfs-dash)
```

## Docker

```bash
docker compose up -d
```

Edit `config.yaml.example` before starting the stack.

## API

- `GET /`
- `GET /api/metrics`
- `GET /api/health/:label`
- `GET /api/health/:label/:pool`
- `GET /health`

### Health Checks

Health checks are designed for monitoring tools like Uptime Kuma. They return `200 OK` if the state is healthy, and `503 Service Unavailable` otherwise.

- **Host Health:** `GET /api/health/:label`
  - Returns `503` if the host is unreachable OR if **no ZFS pools are detected**.
  - Returns `503` if any pool is degraded or faulted.
  - Returns `503` if any pool exceeds `max_usage_percent`.
- **Pool Health:** `GET /api/health/:label/:pool`
  - Returns `503` if the requested pool is **not found** or not healthy.
  - Returns `503` if the pool exceeds `max_usage_percent`.

Examples:

- `GET /api/health/node-1`
- `GET /api/health/node-1/tank`
