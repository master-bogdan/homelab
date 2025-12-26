# Ephemeral Notes API

Ephemeral notes service backed by Redis. Notes are created with a short TTL and
are consumed on first read.

## Features

- REST API for creating, reading, and deleting notes
- Redis storage with 15-minute TTL
- Read-once semantics (notes are deleted on fetch)
- Prometheus metrics endpoint
- Swagger docs
- Simple per-IP rate limiting

## Requirements

- Go 1.22+
- Redis (local via Docker is fine)

## Configuration

Environment variables:

- `SERVER_HOST` (required)
- `SERVER_PORT` (required)
- `REDIS_ADDR` (required, e.g. `localhost:6379`)
- `REDIS_PASSWORD` (optional)

## Run locally

```bash
cd apps/ephermal-notes-api
make run
```

If you prefer manual setup:

```bash
docker-compose -f docker-compose.dev.yaml up -d
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080
export REDIS_ADDR=localhost:6379
go run ./...
```

## API

Base path: `/api/v1`

- `GET /notes/{id}` - fetch and consume a note
- `POST /notes` - create a note
- `DELETE /notes/{id}` - delete a note
- `GET /health/readyz` - readiness check
- `GET /health/healthz` - liveness check
- `GET /metrics` - Prometheus metrics
- `GET /swagger/` - Swagger UI

### Create a note

```bash
curl -sS -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -d '{"message":"hello"}'
```

### Read a note (consumes it)

```bash
curl -sS http://localhost:8080/api/v1/notes/<id>
```

## Notes behavior

- Notes expire after 15 minutes.
- Reading a note deletes it (`GET` is destructive).
- Rate limiting is 5 requests per 10 seconds per IP.

## Development

```bash
make build
make test
make lint
make test-local
```
