# AGENTS.md

## Project

Ephemeral Notes API is a Go service for short-lived notes backed by Redis. Notes have a TTL and are deleted on first read.

## Architecture

- Keep HTTP behavior compatible with the documented `/api/v1` endpoints.
- Preserve read-once note semantics: fetching a note consumes/deletes it.
- Redis is the storage dependency; do not introduce another persistence layer without explicit direction.
- Keep metrics, Swagger, health checks, and per-IP rate limiting intact.
- Keep environment configuration aligned with the README and local compose files.

## Commands

- Run locally with Redis: `make run`
- Build: `make build`
- Unit tests: `make test`
- Local Redis-backed test flow: `make test-local`
- Lint: `make lint`

## Rules

- Do not commit secrets from `.env`.
- Do not change note TTL, rate limits, or destructive-read behavior unless explicitly requested.
- Prefer small, focused Go package changes over broad rewrites.
- Keep Swagger/API docs synchronized when handlers or request/response contracts change.
- Use `gofmt`/`go test` standards for Go changes.

## Verification

For code changes, run at least:

```bash
make test
```

For Redis-dependent behavior, run:

```bash
make test-local
```

For broader changes, also run:

```bash
make lint
make build
```
