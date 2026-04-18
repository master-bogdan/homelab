# AGENTS.md

## Project

Estimate Room API is a Go backend for room estimation workflows. It uses PostgreSQL, Redis, migrations, OAuth/session cookies, Swagger docs, and a WebSocket endpoint.

## Architecture

- Application entrypoint is `cmd/server`.
- Configuration lives in `config/` and is loaded from environment variables, with `.env` support for local development.
- Database schema changes belong in `migrations/`.
- Keep migration files forward/backward paired and compatible with `make migrate-up` and `make migrate-down`.
- Swagger/OpenAPI output is generated through the existing docs tooling.
- Preserve cookie-based auth/session behavior and WebSocket endpoint behavior unless explicitly requested.

## Commands

- Run API: `make run`
- Live reload: `make dev`
- Start local dependencies: `make compose-up`
- Stop local dependencies: `make compose-down`
- Run with dependencies: `make run-compose`
- Test: `make test`
- Lint: `make lint`
- Build: `make build`
- Apply migrations: `make migrate-up`
- Roll back migrations: `make migrate-down`
- Create migration: `make migrate-create name=create_rooms_table`
- Generate Swagger docs: `make swagger-generate`

## Rules

- Do not commit secrets from `.env`.
- Keep `.env.example`, docs, and config validation aligned when environment variables change.
- Do not edit generated Swagger files by hand; update annotations/source and run `make swagger-generate`.
- Do not change auth cookie behavior, OAuth continuation semantics, or WebSocket auth without explicit direction.
- For database changes, include migrations and update tests/docs that rely on schema shape.
- Prefer package-local changes that follow existing patterns in `docs/ARCHITECTURE.md`, `docs/DEVELOPMENT.md`, and `docs/PATTERNS.md`.

## Verification

For backend code changes, run:

```bash
make test
make lint
```

For schema/API contract changes, also run:

```bash
make migrate-up
make swagger-generate
```

For broader changes, run:

```bash
make build
```
