# Frontend Conventions

## Domain Structure

The app is domain-first:

```text
src/
  app/
    pages/
      DashboardPage/
        DashboardPage.tsx
        DashboardPage.styles.ts
        index.ts
      LoginPage/
        LoginPage.tsx
        index.ts
    router/
      AppRouter.tsx
      router.tsx
      routePaths.ts
      routes/
        authRoutes.tsx
        dashboardRoutes.tsx
        index.ts
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
    <module>/
      api/
      store/
      components/
      constants/
      types/
      utils/
      hooks/
      index.ts
  shared/
    api/
    ws/
    store/
    components/
    constants/
    types/
    utils/
    hooks/
    index.ts
```

`app/pages` contains route-entry composers. Organize pages as one folder per page, with page-owned styles beside that page.

Do not group route-entry pages by module under `app/pages/auth`, `app/pages/dashboard`, or similar folders. A page is a composer and can import from any module it needs.

Do not add module `pages` folders for route screens.

`app/router` owns route constants and route composition. Keep `AppRouter.tsx`, `router.tsx`, and `routePaths.ts` at the `app/router` root. Route-array files belong under `app/router/routes`.

Modules do not export route arrays.

## Modules

Use module folders for domain ownership:

- `api`: RTK Query endpoint injection or module API boundary code.
- `store`: `slice.ts`, `selectors.ts`, `types.ts`, `thunks.ts`, and `index.ts`.
- `components`: module-specific reusable UI.
- `constants`: module-specific runtime constants.
- `types`: module-owned typings.
- `utils`: module-specific helpers.
- `hooks`: reusable domain/page orchestration hooks.
- `index.ts`: public module API only.

Do not export every internal utility by default. Export only what outside code intentionally consumes.

## Constants

Use named constants for repeated or behavior-driving string values. Do not compare page, hook, or component state against magic strings inline.

- Constant objects use PascalCase plural names, such as `DashboardLoadStatuses`, `ResetPasswordPageStates`, or `AuthBackToSignInLinkPlacements`.
- Constant object keys are uppercase, such as `LOADING`, `INVALID`, or `CENTERED`.
- Constant object values may remain API/UI string values, such as `'loading'`, `'invalid'`, or `'centered'`.
- Derive union types from the constant object instead of repeating literal unions.
- Module-wide domain constants live in `modules/<module>/constants`.
- Component-only constants may live beside the component when they are not meaningful outside that component.
- Export constants through a module or component public API only when outside callers need them.

## Shared

Use `shared` only for generic, cross-domain code. Do not move domain models or business rules to shared.

Generic examples:

- UI primitives in `shared/components`
- base API client in `shared/api`
- WebSocket client and transport types in `shared/ws`
- theme tokens in `shared/constants`
- app-typed Redux hooks in `shared/hooks`
- typed thunk factory in `shared/store`
- API, theme, and app Redux types in `shared/types`
- date formatting and theme creation in `shared/utils`

Env-backed app config belongs in root `config/appConfig.ts`.

## State And Requests

- Use RTK Query for endpoint-driven server state.
- Use thunks for multi-step business workflows.
- Use slice state for durable domain/UI state that RTK Query does not own.
- Keep React Hook Form state out of Redux unless another app surface needs it.

See [State Management](./architecture/state-management.md) for locations and examples.

## Styling

- MUI theme and `sx` are the primary styling path.
- SCSS Modules are allowed for local page/layout styling.
- `styled()` is limited to reusable wrappers such as shared component building blocks.

## Verification

For structural or request/state refactors, run:

```bash
npm run typecheck
npm test
```

For broader UI or routing changes, run:

```bash
npm run e2e
```
