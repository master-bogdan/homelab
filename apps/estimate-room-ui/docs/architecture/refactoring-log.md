# Refactoring Log

## Moved Files

- Kept app router wiring under `src/app/router`.
- Moved `NotFoundPage` from `src/app/router/NotFoundPage` to `src/app/pages/NotFoundPage`.
- Moved route-entry pages from modules to `src/app/pages`.
- Moved route-entry pages from grouped `app/pages/{auth,dashboard,history,profile,rooms,settings,teams}` folders into one folder per page.
- Moved module route arrays from `src/modules/*/routes.tsx` to `src/app/router/routes/*Routes.tsx`.
- Moved auth page subcomponents from `src/modules/auth/pages/*/components` to `src/modules/auth/components`.
- Moved auth page hooks from `src/modules/auth/pages/*/hooks` to `src/modules/auth/hooks`.
- Moved auth page form types from `src/modules/auth/pages/*/types` to `src/modules/auth/types`.
- Moved auth reset-link utility from `src/modules/auth/pages/ResetPasswordPage/utils` to `src/modules/auth/utils`.
- Moved RTK Query files from `modules/auth/store/authService.ts` and `modules/dashboard/store/dashboardService.ts` into module `api` folders.
- Moved mock module API boundary files from `services` folders into module `api` folders.
- Moved app-typed Redux hooks to `shared/hooks`.
- Moved typed thunk factory to `shared/store`.
- Moved shared UI primitives from `shared/ui` to `shared/components`.
- Restored all env-backed app config under root `src/config/appConfig.ts`.
- Moved theme creation to `shared/utils`, theme tokens to `shared/constants`, and theme types to `shared/types`.
- Restored WebSocket client and transport types under `shared/ws`.
- Moved domain models out of `shared/types/models` into owning module `types` folders.
- Moved `DashboardPage.test.tsx` to `app/pages/DashboardPage/__tests__`.
- Moved `AuthForms.test.tsx` to `app/pages/__tests__`.
- Moved page-owned `RegisterPage` card styling to `app/pages/RegisterPage/RegisterPage.styles.ts`.

## Renamed Files

- Restored `app/store/store.ts` as the store creation file.
- `authStore.ts` -> `store/types.ts`
- `authSlice.ts` -> `store/slice.ts`
- `authSelectors.ts` -> `store/selectors.ts`
- `authThunks.ts` -> `store/thunks.ts`
- `dashboardStore.ts` -> `store/types.ts`
- `dashboardSlice.ts` -> `store/slice.ts`
- `dashboardSelectors.ts` -> `store/selectors.ts`
- `dashboardThunks.ts` -> `store/thunks.ts`
- `systemStore.ts` -> `store/types.ts`
- `systemSlice.ts` -> `store/slice.ts`
- `systemSelectors.ts` -> `store/selectors.ts`

## Created Files

- `docs/architecture/project-structure.md`
- `docs/architecture/state-management.md`
- `docs/architecture/routing.md`
- `docs/architecture/refactoring-log.md`
- `app/guards/index.ts`
- `shared/types/store.ts`
- `shared/ws/index.ts`
- `shared/hooks/useAppDispatch.ts`
- `shared/hooks/useAppSelector.ts`
- app page barrel files under each `app/pages/<PageName>` folder
- app route-array files under `app/router/routes`
- module `api/index.ts` files
- module empty `constants`, `hooks`, `store`, and `components` indexes where needed to normalize structure
- `shared/index.ts`

## Removed Files

- Removed stale module `pages` barrels after route-entry pages moved to `app/pages`.
- Removed stale `shared/theme` and `shared/ui` barrels after content moved to target folders.
- Removed domain model exports from `shared/types`.

## Unresolved Compromises

- `system` still owns global dialog state for dashboard dialogs, including a dashboard payload type. This preserves current behavior but keeps one cross-module UI orchestration dependency.
- Some module components still import app route path constants for link rendering. A future cleanup can pass route targets from pages when those components need stricter independence.
- Mock API boundary files remain promise/local-data based where backend contracts are not yet implemented. They now live in module `api` folders and can be converted to RTK Query endpoint injection when real endpoints exist.
