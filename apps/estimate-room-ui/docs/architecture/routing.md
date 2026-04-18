# Routing

Routing is app-level wiring. Modules do not own route arrays or route-entry page files.

## Pages

Route composers live in `src/app/pages`.

Pages are organized by page folder:

```text
app/pages/
  LoginPage/
    LoginPage.tsx
    index.ts
  DashboardPage/
    DashboardPage.tsx
    DashboardPage.styles.ts
    index.ts
  NewRoomPage/
    NewRoomPage.tsx
    NewRoomPage.module.scss
    index.ts
```

Do not group route-entry pages by module under folders such as `app/pages/auth` or `app/pages/dashboard`. A page is a composer and can pull from any module boundary needed for that route.

Pages may import:

- module hooks, components, store exports, and utils
- shared components/hooks/utils
- app route paths and navigation helpers
- app layouts or guards when needed

Pages should stay thin. They compose module pieces, connect router params, and trigger navigation or redirects.

Examples:

- `app/pages/DashboardPage/DashboardPage.tsx` composes dashboard cards and dashboard hooks.
- `app/pages/JoinRoomPage/JoinRoomPage.tsx` connects route token params to the dashboard join-room thunk.
- `app/pages/LoginPage/LoginPage.tsx` composes auth module forms, cards, and auth hooks.

## Routes

Routes live in `src/app/router`.

```text
app/router/
  AppRouter.tsx
  router.tsx
  routePaths.ts
  routes/
    authRoutes.tsx
    dashboardRoutes.tsx
    historyRoutes.tsx
    profileRoutes.tsx
    roomsRoutes.tsx
    settingsRoutes.tsx
    teamsRoutes.tsx
    index.ts
  index.ts
```

- `AppRouter.tsx`: composer that renders `RouterProvider`.
- `router.tsx`: top-level `createBrowserRouter` composition.
- `routePaths.ts`: route path constants and path builders.
- `routes/*Routes.tsx`: route arrays grouped under one routes folder.

Route paths were preserved during the refactor.

## Layouts

Layouts live in `src/app/layouts`.

- `AuthLayout` handles auth route shell and authenticated-user redirect behavior.
- `DashboardLayout` handles authenticated dashboard shell, session resolution, sidebar/header, and dashboard dialogs.

Layouts may read module selectors and RTK Query hooks when they are app-shell concerns.

## Guards

Reusable route guards belong in `src/app/guards`. The folder is present for future guard extraction. Current auth gating remains in layouts because it is tightly coupled to the existing shell rendering.
