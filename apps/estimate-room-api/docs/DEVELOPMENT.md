# Local Development

## Prerequisites

- Go `1.25.x`
- Docker and Docker Compose
- PostgreSQL
- Redis
- `air` for live reload if you want automatic restarts

The easiest local path is to run PostgreSQL and Redis through Docker Compose and run the API with the local Go toolchain.

## Environment Setup

Copy `.env.example` to `.env` and fill in values.
The API loads `.env` automatically on startup if the file is present.

Minimum variables:

- `HOST`
- `PORT`
- `FRONTEND_BASE_URL`
- `DATABASE_URL`
- `REDIS_URL`
- `PASETO_SYMMETRIC_KEY`
- `ISSUER`

Optional but recommended:

- `WS_ALLOWED_ORIGINS`
- `HTTP_RATE_LIMIT_PER_MINUTE`
- `WS_RATE_LIMIT_PER_MINUTE`
- GitHub OAuth variables if GitHub login is required
- Email SMTP variables if local password reset emails should be delivered

`FRONTEND_BASE_URL` is used by `/api/v1/oauth2/authorize` to redirect unauthenticated users to the frontend login page with a `continue` URL.

The checked-in local defaults are aligned to the dev compose stack:

- API: `http://localhost:8080`
- Frontend: `http://localhost:5173`
- PostgreSQL: `postgres://postgres:password@localhost:5432/estimate_room?sslmode=disable`
- Redis: `redis://localhost:6379/0`
- SMTP: `localhost:1025`
- Mail inbox UI: `http://localhost:8025`

## Start Dependencies

```bash
make compose-up
```

This starts the local PostgreSQL and Redis stack from `docker-compose.dev.yaml`, plus:

- pgAdmin: `http://localhost:5050`
- RedisInsight: `http://localhost:5540`
- Mailpit: `http://localhost:8025`

Default pgAdmin login:

- Email: `admin@estimate-room.dev`
- Password: `admin`

When connecting from pgAdmin to the database container, use:

- Host: `postgres`
- Port: `5432`
- Username: `postgres`
- Password: `password`
- Database: `estimate_room`

When connecting from RedisInsight to the Redis container, use:

- Host: `redis`
- Port: `6379`

When connecting the API to the local mail inbox, use:

- SMTP host: `localhost`
- SMTP port: `1025`
- Mailpit UI: `http://localhost:8025`

## Run the API

```bash
make run
```

With the default local `.env`, `IS_AUTO_MIGRATIONS=true`, so the API applies pending migrations at boot.
If you disable auto-migrations, run `make migrate-up` manually before starting the API.

## Live Reload

Install `air` once:

```bash
go install github.com/air-verse/air@latest
```

Run the API with automatic restart on code or `.env` changes:

```bash
make dev
```

If you also want the Docker dependencies started first:

```bash
make dev-compose
```

## Useful Commands

Build:

```bash
make build
```

Run tests:

```bash
make test
```

Lint:

```bash
make lint
```

Apply migrations:

```bash
make migrate-up
```

Rollback migrations:

```bash
make migrate-down
```

## Migration Notes

- Migration files live in `migrations/`.
- The initial schema is already checked in as `1767825448_initial.*.sql`.
- Auto-migrations only run when `IS_AUTO_MIGRATIONS=true`.
- Manual migration commands use `DATABASE_URL` from your `.env`.
- The schema does not seed OAuth clients. If you need the browser OAuth2 authorization flow, insert at least one row into `oauth2_clients`.

Generate Swagger docs:

```bash
make swagger-generate
```

Stop local infrastructure:

```bash
make compose-down
```

## Important Endpoints

HTTP API base:

```text
/api/v1
```

Health:

- `/api/v1/health/healthz`
- `/api/v1/health/readyz`

Auth:

- `/api/v1/auth/login`
- `/api/v1/auth/register`
- `/api/v1/auth/session`
- `/api/v1/auth/logout`
- `/api/v1/auth/forgot-password`
- `/api/v1/auth/reset-password/validate`
- `/api/v1/auth/reset-password`
- `/api/v1/auth/github/login`
- `/api/v1/auth/github/callback`

Swagger UI:

- `/swagger/index.html`

Metrics:

- `/metrics`

## Testing Notes

The test suite expects a PostgreSQL database.

Use one of:

- `TEST_DATABASE_URL`
- `DATABASE_URL`

The tests run migrations automatically before opening the DB connection.

## WebSocket Notes

- WebSocket endpoint: `/api/v1/ws`
- Registered users authenticate with the access token
- Guests authenticate with the room guest cookie
- Query-string access tokens are intentionally rejected

## Known Local Gaps

- GitHub OAuth requires real client credentials and redirect URLs.
- Redis pub/sub is required for multi-instance websocket fan-out, but single-process tests can run with an in-memory pub/sub stub.
- The checked-in OpenAPI output is partial and should be treated as incomplete until broader endpoint annotation work is done.
