# Engineering Patterns

## Architectural Pattern

The codebase follows a module-oriented service architecture inside a single binary.

Each business area lives under `internal/modules/<name>` and usually contains:

- `*_module.go`: wiring and route registration
- `*_controller.go`: transport concerns
- `*_service.go`: business logic and authorization
- `repositories/`: persistence access
- `models/`: persistence/domain models
- `dto/`: transport contracts
- `tests/`: package-level behavioral tests

## Controller Pattern

Controllers should:

- parse request input
- call auth extraction when needed
- validate DTOs early
- delegate business rules to services
- translate domain/application errors into HTTP responses

Controllers should not:

- contain authorization policy beyond simple identity extraction
- perform multi-step persistence work
- duplicate business invariants already owned by services

## Service Pattern

Services are the primary home for:

- authorization decisions
- domain invariants
- transactional orchestration
- cross-repository workflows
- side effects such as rewards, expiry touches, and broadcasts

When behavior spans multiple aggregates, prefer a service method over controller-level orchestration.

## Repository Pattern

Repositories should stay narrow and persistence-focused:

- fetch by identity or relation
- create/update/delete data
- keep SQL close to the table it owns

Avoid pushing product policy into repositories. If the rule reads like business language, keep it in the service layer.

## Authorization Pattern

Authorization is enforced in services, not only in controllers.

Use the service layer to answer questions like:

- Is this user a room admin?
- Is this user a member of the attached team?
- Is this participant eligible to vote?
- Can this actor revoke this invite?

This keeps transport-specific entrypoints from drifting apart.

## Error Handling Pattern

Use typed application errors from `internal/pkg/apperrors`.

Expected categories:

- `ErrBadRequest`
- `ErrUnauthorized`
- `ErrForbidden`
- `ErrNotFound`
- `ErrConflict`
- `ErrInternal`

Transport adapters convert these into `application/problem+json`.

Rules:

- return precise domain errors from services
- keep controller error mapping consistent
- do not leak raw DB or implementation details to clients

## Transaction Pattern

Use `bun` transactions for multi-step state changes that must remain atomic.

Examples:

- room creation plus initial admin participant
- room creation plus invite fan-out
- expiry state transitions
- terminal reward application

If a workflow can partially succeed safely, document why. Otherwise, wrap it in a transaction.

## Realtime Pattern

The websocket layer is split in two parts:

- `ws.Service`: connection lifecycle, presence, outbound fan-out
- `roomsGateway`: room-specific realtime behaviors and event semantics

Guidelines:

- keep transport concerns in `ws`
- keep room business semantics in `roomsGateway`
- publish normalized event envelopes
- send snapshots on join/reconnect
- keep websocket payloads explicit and version-safe

## State Mutation Pattern

Any successful room mutation should touch `last_activity_at`.

This includes:

- room updates
- invite joins
- task changes
- vote lifecycle events
- websocket room joins/leaves that affect presence

That keeps expiry behavior deterministic.

## Observability Pattern

Every feature should produce:

- request logs for the HTTP boundary
- focused service or realtime logs for important lifecycle changes
- metrics for high-value operational counters

Prefer structured logs with stable keys over interpolated strings.

## Testing Pattern

Use three layers of tests:

- unit-like service tests for domain rules
- controller/module tests for HTTP behavior and auth mapping
- end-to-end tests for full user flows

Guidelines:

- seed only the tables required for the scenario
- assert business outcomes, not implementation details
- test both allowed and denied paths for authorization-sensitive endpoints
- cover realtime flows with explicit event assertions

## Extension Rules

When adding a feature:

1. add or extend a module instead of creating cross-cutting ad hoc code
2. keep auth and authorization in service methods
3. update docs if the runtime behavior changes
4. add at least one regression test on the critical path
5. update `docs/DEV_PLAN.md` if the feature changes delivery status
