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

## Config

```yaml
addr: ":8054"
refresh: 300
debug: false
trusted_proxies: [] # List of proxy IPs (e.g., ["127.0.0.1", "10.0.0.1"])

endpoints:
  - url: "http://host1:9134/metrics"
    label: "node-1"
    location: "Singapore"
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
- **Pool Health:** `GET /api/health/:label/:pool`
  - Returns `503` if the requested pool is **not found** or not healthy.

Examples:

- `GET /api/health/node-1`
- `GET /api/health/node-1/tank`
