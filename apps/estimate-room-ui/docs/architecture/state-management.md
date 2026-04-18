# State Management

Redux Toolkit is the standard state layer. RTK Query owns server state. Thunks own multi-step workflows. Slices own durable domain/UI state that RTK Query does not already model.

## App Store

App-wide Redux wiring lives in `src/app/store`:

- `store.ts`: `configureStore`.
- `index.ts`: barrel exports only.
- `rootReducer.ts`: combines RTK Query and module reducers.
- `middleware.ts`: custom middleware declarations, currently RTK Query middleware.

App Redux types live in `src/shared/types/store.ts`:

- `RootState`
- `AppDispatch`

App-typed Redux hooks live in `src/shared/hooks`:

- `useAppDispatch.ts`
- `useAppSelector.ts`

The typed thunk factory lives in `src/shared/store/createAppAsyncThunk.ts`.

The root reducer currently combines:

- `shared/api` RTK Query reducer
- `auth` reducer
- `dashboard` reducer
- `system` reducer

## Module Store

Each module has a `store` folder with the same standard files:

- `slice.ts`
- `selectors.ts`
- `types.ts`
- `thunks.ts`
- `index.ts`

Modules without Redux state still keep empty files so the module shape stays predictable.

## RTK Query

RTK Query endpoint injection lives in module `api` folders.

Examples:

- `modules/auth/api/authApi.ts`: login, register, logout, session, password reset, OAuth token exchange.
- `modules/dashboard/api/dashboardApi.ts`: dashboard sessions, teams, ledger, room preview, invite preview, create room, accept invite.

Use RTK Query for:

- simple CRUD
- list/detail fetching
- create/update/delete requests
- endpoint-driven server state
- caching and invalidation

## Thunks

Thunks live in `modules/<name>/store/thunks.ts`.

Use `createAppAsyncThunk` from `shared/store` for app-typed thunk creation.

Current examples:

- `completeOAuthCallback`: validates OAuth state, exchanges authorization code, fetches session, updates auth state.
- `fetchDashboardPage`: fetches sessions, teams, ledger, and active room, then composes dashboard page state.
- `submitJoinRoom`: parses invite token, previews invite, rejects team invites, accepts room invite.
- `submitCreateRoom`: creates a room and refreshes dashboard state.

Use thunks for:

- multiple dependent API calls
- branching based on state or API results
- workflows dispatching several actions
- orchestration across module boundaries

## Slice State

Use slice state only for durable domain state or composed UI state that RTK Query does not own.

Current examples:

- `auth`: current session user, auth status, OAuth callback workflow state.
- `dashboard`: composed dashboard page data and dialog workflow state.
- `system`: global dialogs, notifications, sidebar state, theme mode.

Do not mirror RTK Query request state in slices unless the state is composed workflow state.
