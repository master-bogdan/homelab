# Project Structure

EstimateRoom UI is organized domain-first. Technical wiring lives in `app`, domain behavior lives in `modules`, and generic building blocks live in `shared`.

## Final Tree

```text
src/
  app/
    App.tsx
    pages/
      DashboardPage/
        DashboardPage.tsx
        DashboardPage.styles.ts
        index.ts
      LoginPage/
        LoginPage.tsx
        index.ts
      NewRoomPage/
        NewRoomPage.tsx
        NewRoomPage.module.scss
        index.ts
    router/
      AppRouter.tsx
      router.tsx
      routePaths.ts
      routes/
        authRoutes.tsx
        dashboardRoutes.tsx
        historyRoutes.tsx
        index.ts
    layouts/
    providers/
    guards/
    store/
      store.ts
      index.ts
      rootReducer.ts
      middleware.ts
  config/
    appConfig.ts
    index.ts
  modules/
    auth/
      api/
      store/
      components/
      constants/
      types/
      utils/
      hooks/
      index.ts
    dashboard/
      api/
      store/
      components/
      constants/
      types/
      utils/
      hooks/
      index.ts
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
    index.ts
```

Every module follows the same internal shape: `api`, `store`, `components`, `constants`, `types`, `utils`, `hooks`, and `index.ts`.

## App

`app` is top-level application wiring only:

- `pages` contains route-entry page composers.
- `router` owns route constants, route arrays, and router composition.
- `layouts` owns shell layout components and auth/dashboard layout decisions.
- `providers` owns app provider composition.
- `guards` is reserved for reusable route guards.
- `store` owns Redux bootstrap.

`app/store/store.ts` owns `configureStore`. `app/store/index.ts` is a barrel only. App Redux types live in `shared/types/store.ts`.

Pages live in `app/pages` because they are route-entry composition points. They connect router params, navigation, redirects, and module pieces, but reusable business logic stays in modules.

Pages are organized as one folder per route-entry page, not by module or domain. Page-owned styles and page-local tests live beside that page:

```text
app/pages/
  RegisterPage/
    RegisterPage.tsx
    RegisterPage.styles.ts
    index.ts
  DashboardPage/
    DashboardPage.tsx
    DashboardPage.styles.ts
    __tests__/
    index.ts
```

Each page may compose anything it needs from multiple modules. For example, a page in `app/pages/TeamDetailsPage` can compose team data, auth/user information, shared UI primitives, and app route helpers without becoming part of the `teams` module.

## Modules

Modules own business/domain code. Current modules are:

- `auth`
- `dashboard`
- `history`
- `profile`
- `rooms`
- `settings`
- `system`
- `teams`

Folder responsibilities:

- `api`: RTK Query endpoints or module API boundary code.
- `store`: Redux selectors, slices, thunks, and store metadata.
- `components`: module-specific reusable UI.
- `constants`: module-specific runtime constants.
- `types`: module-owned typings.
- `utils`: module-specific helpers.
- `hooks`: reusable module/page orchestration hooks.
- `index.ts`: public API only.

## Shared

`shared` contains only generic reusable code:

- `api`: base RTK Query API, base query, token/session lifecycle.
- `ws`: WebSocket client and related transport types.
- `components`: generic UI primitives.
- `constants`: generic constants such as theme tokens.
- `hooks`: generic hooks and app-typed Redux hooks.
- `types`: generic API, theme, and app Redux types.
- `utils`: generic formatters and theme creation.
- `store`: generic store utilities, including the typed thunk factory.

Domain models were moved out of `shared/types` into their owning modules.

Env-backed runtime config lives in root `config/appConfig.ts`.

## Public APIs

Outside code should import module behavior from the module public API or from a deliberate public subfolder API such as `modules/auth/components`.

Do not broadly export internal constants, utility functions, or transport types by default. Export only code that is intentionally consumed outside the module.
