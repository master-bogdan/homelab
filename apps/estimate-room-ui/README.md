# estimate-room-ui

Frontend for EstimateRoom room planning, dashboard workflows, history review, team pages, profile/settings, and backend API/WebSocket integration.

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
    pages/
      DashboardPage/
      LoginPage/
      NewRoomPage/
    router/
      AppRouter.tsx
      router.tsx
      routePaths.ts
      routes/
    layouts/
    providers/
    guards/
    store/
      store.ts
      index.ts
  config/
    appConfig.ts
    index.ts
  modules/
    auth/
    dashboard/
    history/
    profile/
    rooms/
    settings/
    system/
    teams/
  shared/
    api/
    ws/
    store/
    components/
    constants/
    types/
    utils/
    hooks/
```

## Conventions

- Vite + React + TypeScript + npm, with no Next.js.
- `app` wires pages, router, layouts, providers, guards, and Redux bootstrap.
- `app/pages` contains route-entry page composers only, organized as one folder per page.
- Page-owned styles live beside the page, for example `app/pages/NewRoomPage/NewRoomPage.module.scss`.
- `modules/*` owns business/domain logic.
- `shared/*` owns generic reusable primitives only.
- RTK Query endpoint injection lives in module `api` folders.
- WebSocket client and transport types live in `shared/ws`.
- Behavior-driving status/state/placement values should use named PascalCase constant objects with uppercase keys, not inline magic strings.
- Redux selectors, slices, thunks, and store metadata live in module `store` folders.
- App-typed Redux hooks live in `shared/hooks`.
- App Redux types live in `shared/types/store.ts`.
- The typed thunk factory lives in `shared/store`.
- Env-backed app config lives in root `config/appConfig.ts`.
- Components stay presentation-focused.
- Hooks orchestrate page and domain behavior.
- Form state stays in React Hook Form instead of Redux.

Architecture docs:

- [`docs/architecture/project-structure.md`](./docs/architecture/project-structure.md)
- [`docs/architecture/state-management.md`](./docs/architecture/state-management.md)
- [`docs/architecture/routing.md`](./docs/architecture/routing.md)
- [`docs/architecture/refactoring-log.md`](./docs/architecture/refactoring-log.md)

Additional engineering notes:

- [`docs/STATE_AND_REQUEST_MANAGEMENT.md`](./docs/STATE_AND_REQUEST_MANAGEMENT.md)
- [`docs/REQUEST_FLOW_DECISIONS.md`](./docs/REQUEST_FLOW_DECISIONS.md)
- [`docs/FRONTEND_CONVENTIONS.md`](./docs/FRONTEND_CONVENTIONS.md)
