# EstimateRoom Architecture

## Purpose

EstimateRoom is a realtime estimation backend for planning poker style sessions. It supports:

- Email/password sign-in through a local auth module that resumes the OAuth2 authorization-code flow
- GitHub sign-in
- Password reset for local accounts
- Teams and invitations
- Rooms with share links and guest access
- Task estimation rounds over WebSockets
- History and lightweight gamification

The backend is implemented as a single Go service with PostgreSQL for persistence and Redis for cross-instance websocket fan-out.

## Runtime Topology

### Core runtime pieces

- HTTP server: `chi` router, middleware stack, JSON/problem+json responses
- Auth: `/auth/*` for browser session and account lifecycle plus `/oauth2/*` for authorization-code + PKCE token issuance
- Realtime: WebSocket transport with Redis-backed pub/sub for broadcast fan-out
- Persistence: PostgreSQL via `bun`
- Background work: room inactivity expiry sweep
- Observability: structured logs, Prometheus-compatible `/metrics`, health probes

### External dependencies

- PostgreSQL
- Redis
- GitHub OAuth endpoints

## High-Level Request Flow

### HTTP flow

1. Request enters global middleware.
2. Request ID, request logging, panic recovery, and IP rate limiting are applied.
3. The request is routed under `/api/v1/...` or to operational endpoints such as `/metrics` and `/swagger/*`.
4. Controllers decode input and map domain/application errors to HTTP responses.
5. Services enforce business rules and authorization.
6. Repositories perform persistence operations.

### WebSocket flow

1. Client connects to `/api/v1/ws`.
2. The backend authenticates either:
   - a registered user via access token, or
   - a guest via the room guest cookie
3. A single active socket per identity is enforced.
4. Room-specific events are handled by the rooms gateway.
5. Outbound broadcasts are published through the local WS service and Redis pub/sub.
6. Reconnects receive a fresh room snapshot.

## Module Map

### `auth`

- Local email/password login
- Local account registration
- Browser session inspection and logout
- Password reset token validation and password reset
- GitHub login bridge for first-party auth

### `oauth2`

- Authorization-code + PKCE flow
- Access token and refresh token issuance
- Refresh token rotation by revoke-and-reissue
- OIDC session persistence for browser auth continuation

### `users`

- Current user lookup via `/users/me`
- User persistence and GitHub profile linking

### `invites`

- Team invites
- Room email invites
- Room share-link invites
- Guest room token issuance and validation

### `teams`

- Team creation
- Team membership reads
- Team invite creation
- Team member removal

### `rooms`

- Room creation and update
- Task CRUD
- Voting round lifecycle
- Final estimate persistence
- Expiry sweep for inactive rooms
- Realtime gateway for room collaboration

### `history`

- Personal session history
- Team session history
- Room summary

### `gamification`

- User stats
- Achievement unlocks
- Room completion rewards
- Realtime reward notifications

### `health`

- Liveness and readiness probes

### `ws`

- WebSocket connection management
- Identity binding
- Presence tracking
- Reconnect snapshots
- Inbound websocket message rate limiting

## Data Model

### Identity and auth

- `users`
- `auth_password_reset_tokens`
- `oauth2_clients`
- `oauth2_oidc_sessions`
- `oauth2_auth_codes`
- `oauth2_refresh_tokens`
- `oauth2_access_tokens`

### Collaboration

- `teams`
- `team_members`
- `rooms`
- `room_participants`
- `tasks`
- `task_rounds`
- `votes`
- `invitations`

### Personalization and progression

- `user_settings`
- `user_stats`
- `user_achievements`
- `user_session_rewards`
- `team_stats`
- `team_achievements`

### Reference data

- `decks`

## Key Domain Rules

- Only registered users can create rooms.
- A room creator is inserted as the room admin participant.
- Team-attached rooms require team membership.
- Only room admins can mutate room/task state.
- Only eligible participants can vote in the active round.
- Only one active task may exist per room.
- Final estimate values must come from the room deck.
- Guests can read only the room they joined through a valid guest token.
- Resetting a password revokes all active browser sessions and tokens for that user.
- Inactive active rooms are expired by the background sweep.

## Realtime Model

### Core incoming events

- `ROOMS_JOIN`
- `ROOMS_TASK_SET_CURRENT`
- `ROOMS_VOTE_CAST`
- `ROOMS_VOTE_REVEAL`
- `ROOMS_ROUND_NEXT`
- `ROOMS_TASK_FINALIZE`

### Core outgoing events

- `ROOMS_SNAPSHOT`
- `ROOMS_PARTICIPANT_JOINED`
- `ROOMS_PARTICIPANT_LEFT`
- `ROOMS_TASK_CURRENT_CHANGED`
- `ROOMS_VOTE_STATUS_CHANGED`
- `ROOMS_VOTES_ALL_CAST`
- `ROOMS_VOTES_REVEALED`
- `ROOMS_ROUND_CHANGED`
- `ROOMS_TASK_FINALIZED`
- `ROOMS_EXPIRED`

## Background Processing

The expiry service runs continuously once the app boots:

- It scans for active rooms whose `last_activity_at` is older than the expiry threshold.
- It marks them expired in a transaction.
- It applies terminal rewards best-effort.
- It emits the `ROOMS_EXPIRED` event.

## Observability

### Logs

- Structured JSON logs
- Request logs with route, status, duration, and request ID
- WebSocket connect, disconnect, and error logs
- Room lifecycle logs

### Metrics

- `/metrics`
- HTTP request totals and latency histogram
- Active and total websocket connections
- Room lifecycle counters for create/finish/expire

### Health

- `/api/v1/health/healthz`
- `/api/v1/health/readyz`

## Security Posture

### Implemented controls

- PKCE for local OAuth2 authorization-code flow
- Signed GitHub OAuth state with expiry
- Access token lookup against persisted token records
- Refresh token rotation by revocation
- HttpOnly session and guest cookies
- SameSite=Lax session cookie
- HTTP IP rate limiting
- WebSocket inbound message rate limiting
- Room/team role checks in service layer

### Still-open product gaps

- No logout/session revocation endpoint yet
- No profile edit or theme API yet
- No committed dashboards/alerts assets yet
- OpenAPI coverage does not yet match the full implemented API surface
- Auto-reveal-on-all-votes is not implemented in the backend yet
- Room-finished realtime event is still missing

## Deployment Notes

- Stateless HTTP nodes are expected.
- PostgreSQL is the source of truth.
- Redis is required for multi-instance websocket broadcast consistency.
- Graceful shutdown stops the expiry loop, shuts down WS, then closes Redis and DB connections.
