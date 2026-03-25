# Local Development

## Prerequisites

- Go `1.25.x`
- Docker and Docker Compose
- PostgreSQL
- Redis

The easiest local path is to run PostgreSQL and Redis through Docker Compose and run the API with the local Go toolchain.

## Environment Setup

Copy `.env.example` to `.env` and fill in values.

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

`FRONTEND_BASE_URL` is used by `/api/v1/oauth2/authorize` to redirect unauthenticated users to the frontend login page with a `continue` URL.

## Start Dependencies

```bash
make compose-up
```

This starts the local PostgreSQL and Redis stack from `docker-compose.dev.yaml`.

## Run the API

```bash
make run
```

If `IS_AUTO_MIGRATIONS=true`, the API applies pending migrations at boot.

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
