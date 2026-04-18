# AGENTS.md

## Project

Estimate Room UI is a Vite + React + TypeScript application for room estimation, dashboard workflows, history, teams, profile/settings, and API/WebSocket integration.

## Architecture

- `src/app` is top-level wiring only: pages, router, layouts, providers, guards, and Redux bootstrap.
- `src/app/pages` uses one folder per route-entry page. Pages are thin composers and may compose from any module.
- `src/app/router` owns `AppRouter.tsx`, `router.tsx`, `routePaths.ts`, and route-array files under `router/routes/`.
- `src/modules/*` owns domain logic, with `api`, `store`, `components`, `constants`, `types`, `utils`, `hooks`, and `index.ts`.
- `src/shared/*` is for generic reusable code only. Do not move domain behavior into shared.
- `src/config/appConfig.ts` is the single env-backed app config file.
- `src/shared/ws` contains WebSocket client and transport types.
- `src/app/store/store.ts` creates the Redux store. `src/app/store/index.ts` is export-only.
- Typed Redux hooks live in `src/shared/hooks`.
- `RootState` and `AppDispatch` live in `src/shared/types/store.ts`.
- Typed thunk factory lives in `src/shared/store`.

## Commands

- Dev server: `npm run dev`
- Build: `npm run build`
- Typecheck: `npm run typecheck`
- Unit tests: `npm test`
- Watch tests: `npm run test:watch`
- Coverage: `npm run coverage`
- Lint: `npm run lint`
- E2E: `npm run e2e`

## Rules

- Do not change route paths unless explicitly requested.
- Do not put route-entry screens inside module `pages` folders.
- Do not group app pages by module; use page folders such as `src/app/pages/DashboardPage/`.
- Do not dump implementation into `index.ts`; index files are public export barrels only.
- Do not use magic strings for behavior-driving state/status/placement values; create PascalCase plural constant objects with uppercase keys.
- Use RTK Query for endpoint-driven server state and thunks for multi-step orchestration.
- Keep module public APIs intentional; do not broadly export internal constants/types/utils.
- Keep page-owned styles beside the page.
- Preserve business behavior during structural refactors.

## Verification

For most code changes, run:

```bash
npm run typecheck
npm test
npm run lint
```

For routing, layout, or browser behavior changes, also run:

```bash
npm run e2e
```

Update architecture docs under `docs/architecture/` when structural boundaries change.
