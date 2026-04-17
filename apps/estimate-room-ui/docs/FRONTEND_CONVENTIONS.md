# Frontend Conventions

## Module Structure

Each module should keep ownership clear:

```text
module/
  components/
  constants/
  hooks/
  pages/
  store/
  types/
  utils/
  routes.tsx
```

Use `routes.tsx`, not names with extra dot notation such as `auth.routes.tsx`.

Route page implementations live under `pages/`. Module root files are for
reusable module API, not page implementation files.

Use this shape for route pages:

```text
module/pages/LoginPage/
  components/
    LoginForm/
      LoginForm.tsx
      index.ts
  hooks/
    useLoginPage.ts
    index.ts
  styles/
    index.ts
  types/
    index.ts
  LoginPage.tsx
  index.ts
```

Page rules:

- Page components own render structure and simple view wiring.
- Page-specific business logic lives beside the page in a focused hook.
- Page-specific UI blocks live under the page `components/` folder.
- Page-specific components use their own folders with an `index.ts`.
- Page-specific hooks live in the page `hooks/` folder.
- Shared hooks live in `module/hooks` only when more than one page or feature uses them.
- Page styles live in the page `styles/` folder.
- Page-specific types live in the page `types/` folder.
- Page-specific constants live in the page `constants/` folder when needed.
- Do not dump loose `use*.ts`, `types.ts`, `styles.ts`, or constants files beside the page component.
- Reusable module components stay in `module/components`.

Route arrays use PascalCase:

```ts
export const AuthRoutes = [...]
export const DashboardRoutes = [...]
```

## Types Folder

`types/` is type-only.

Allowed:

- `type`
- `interface`
- type-only imports

Not allowed:

- `const`
- runtime values
- `as const`
- functions
- request logic

Runtime constants live in `constants/`.

Type files can derive from constants using type-only imports:

```ts
import type { AuthStates } from '../constants';

export type AuthStatus = (typeof AuthStates)[keyof typeof AuthStates];
```

## API Type Names

Do not use `Dto` naming in the UI.

Use:

- `ApiRequest`
- `ApiResponse`
- `Payload`
- `Result`
- `Model`

Examples:

```ts
SessionApiResponse
OAuthTokenApiResponse
DashboardCreateRoomApiResponse
DashboardCreateRoomResult
```

Transport response types can live in `types/api.ts`.

Domain models should live in `types/models.ts`.

Form values should live in `types/forms.ts`.

## Constants

Runtime constants use PascalCase for exported objects and uppercase keys:

```ts
export const AppRoutes = {
  DASHBOARD: '/dashboard',
  LOGIN: '/login'
} as const;

export const AuthStates = {
  AUTHENTICATED: 'authenticated',
  UNAUTHENTICATED: 'unauthenticated',
  UNKNOWN: 'unknown'
} as const;
```

Do not export lowercase constant bags such as:

```ts
appRoutes
authStatuses
appConfig
```

Use:

```ts
AppRoutes
AuthStates
AppConfig
```

`AppConfig` keys are uppercase.

## Component Props

Props are owned by the component file.

Do not export props from module barrels.

Allowed:

```ts
export interface AuthCardProps {
  readonly children: ReactNode;
}
```

Not allowed in `components/index.ts`:

```ts
export type { AuthCardProps } from './AuthCard';
```

## Components

Components should:

- render UI
- call RTK hooks for simple endpoint requests
- dispatch high-level domain thunks when the flow is multi-step
- read selectors
- pass explicit event handlers to child components

Components should not:

- call `.unwrap()`
- contain API orchestration
- duplicate RTK Query loading state
- own cross-module request coordination
- use nested ternaries for meaningful UI branches

Extract meaningful local variables for complex conditions:

```ts
const isResolvingSession =
  shouldFetchSession &&
  (sessionQuery.isUninitialized || sessionQuery.isLoading || sessionQuery.isFetching);

if (isResolvingSession) {
  return <AppPageState isLoading title="Loading session" />;
}
```

## Forms

Use React Hook Form for form state.

Keep submitted form values in form or local page state unless another part of the app needs that state.

Do not move form-only state into Redux.

## Barrels

Barrels export public module API only.

Good exports:

- pages
- routes
- public hooks
- public selectors
- public domain thunks
- RTK hooks intended for module consumers

Avoid exporting:

- component prop types
- private helpers
- transport-only types unless another module truly needs them

## Styling

Separate styles when a page or component has meaningful styling logic.

Preferred local style file:

```text
ComponentName.tsx
ComponentName.styles.ts
index.ts
```

Avoid embedding large style objects in component render files.

## Verification

For request/state refactors, run:

```bash
npm run typecheck
npm test -- AuthLayout DashboardLayout AuthForms useOAuthCallbackPage dashboardThunks
```

For broader UI changes, run:

```bash
npm test
npm run e2e
```
