# zfs-dash

Minimal dashboard for ZFS pool metrics from `pdf/zfs_exporter`.

## Run

```bash
cp config.yaml.example config.yaml
go run . serve --config config.yaml
```

Open `http://localhost:8054`.

## Config

```yaml
addr: ":8054"
refresh: 30

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
- `GET /health`
