# EstimateRoom UI Development Docs

This folder defines frontend engineering rules for EstimateRoom UI.

Read these before changing request flows, Redux state, module structure, or shared conventions:

- [State And Request Management](./STATE_AND_REQUEST_MANAGEMENT.md)
- [Request Flow Decisions](./REQUEST_FLOW_DECISIONS.md)
- [Frontend Conventions](./FRONTEND_CONVENTIONS.md)

These rules are meant to keep request ownership explicit:

- RTK Query owns simple API requests.
- Thunks own multi-step domain workflows.
- Components own presentation and user-facing error display.
- Redux slices own domain state, not copied request state.

