# estimate-room-ui

Frontend scaffold for room estimation workflows, history review, team pages, profile/settings, and future Go backend API plus WebSocket integration.

## Install

```bash
npm install
```

## Run

```bash
npm run dev
```

## Test

```bash
npm test
npm run test:watch
npm run coverage
npm run typecheck
```

## Folder Structure

```text
src/
  app/
    layouts/
    providers/
    router/
    store/
  modules/
    auth/
    dashboard/
    history/
    profile/
    rooms/
    settings/
    teams/
  shared/
    api/
    config/
    constants/
    hooks/
    types/
    ui/
    utils/
    ws/
  assets/
  styles/
  test/
  theme/
```

## Conventions

- Vite + React + TypeScript + npm, with no Next.js.
- One repository-level `tsconfig.json` covers both app code and `vite.config.ts`.
- MUI theme and `sx` are the primary styling path.
- SCSS Modules are reserved for local page or layout styling.
- `styled()` is limited to reusable wrappers such as shared UI building blocks.
- Redux Toolkit is kept deliberate: `auth` and `ui` are the initial global slices.
- Components stay presentation-focused.
- Hooks orchestrate effects and page logic.
- Services own API and WebSocket boundaries.
- Selectors live outside slices.
- Form state stays in React Hook Form instead of Redux.
- Path aliases are available through `@`, including `@/app`, `@/modules`, `@/shared`, and `@/theme`.

Detailed engineering rules live in [`docs/`](./docs/):

- [`docs/STATE_AND_REQUEST_MANAGEMENT.md`](./docs/STATE_AND_REQUEST_MANAGEMENT.md)
- [`docs/REQUEST_FLOW_DECISIONS.md`](./docs/REQUEST_FLOW_DECISIONS.md)
- [`docs/FRONTEND_CONVENTIONS.md`](./docs/FRONTEND_CONVENTIONS.md)

## Notes

- Authentication is scaffolded, but there is no fake login implementation.
- Shared API and WebSocket clients are ready for Go backend integration.
- `src/modules/rooms/NewRoomPage.tsx` demonstrates the preferred form pattern for future feature pages.
