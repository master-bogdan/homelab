# State And Request Management

## Core Rule

Use the simplest owner that can correctly handle the behavior.

- One endpoint request: use an RTK Query hook.
- Multi-step workflow: use a thunk.
- Presentation and error display: keep it in the component or page hook.
- Domain state: keep it in a Redux slice.

Do not wrap simple RTK Query endpoints in thunks just to start a request.

## RTK Query Hooks

Use RTK Query hooks directly for simple CRUD and single endpoint calls.

Examples:

- login
- register
- logout
- forgot password
- reset password
- validate reset token
- fetch current session
- create room, when it is only the create request

Allowed in components and page hooks:

```ts
const [logout, logoutState] = useLogoutMutation();
const sessionQuery = useFetchSessionQuery();
```

Components may use RTK state directly:

- `data`
- `error`
- `isLoading`
- `isFetching`
- `isError`
- `isSuccess`

Components must not call `.unwrap()` for RTK requests.

Components must not add local duplicated request state such as:

- `isLoggingOut`
- `isSubmitting`
- `isLoadingTeams`

when RTK Query already exposes the same state.

## Page Hooks

Page hooks are allowed only when they own page orchestration that would make the
page component hard to read.

Good page hook responsibilities:

- React Hook Form setup for a large form
- route/search-param interpretation
- page-only state machines
- mapping transport errors into field or page errors
- coordinating an RTK request with redirect behavior
- dispatching a high-level thunk for a multi-step workflow

Do not create a page hook only to rename an RTK Query hook or wrap one endpoint.
If a one-endpoint feature still has coherent stateful behavior around it, extract
that behavior as a focused flow hook instead of a whole page hook.

Bad page hook:

```ts
const [forgotPassword, forgotPasswordState] = useForgotPasswordMutation();
```

when it just renames the mutation and returns the same request state.

Good focused flow hook:

```ts
const { errorMessage, isSubmitted, resend, submit } = useForgotPasswordFlow();
```

when the hook owns mutation submission, submitted email state, resend behavior,
and API error normalization while the page keeps form setup and rendering.

## Thunks

Use thunks for domain workflows that coordinate more than one thing.

Use a thunk when the action needs:

- multiple API calls
- sequential branching
- derived request decisions
- cross-slice updates
- composed page state
- one high-level domain intent from the UI

Examples:

- OAuth callback: validate state, exchange code, fetch session, update auth state
- dashboard page load: fetch sessions, teams, ledger, active room, then compose view state
- join room: parse invite token, preview invite, reject team invite, accept room invite

Components may dispatch thunks directly. That is normal Redux usage.

The important rule is that components dispatch the highest-level intent they know about:

```ts
dispatch(fetchDashboardPage());
dispatch(joinRoomFromInvite(code));
```

Avoid transport-only thunk names and behavior such as:

```ts
dispatch(startLogoutRequest());
dispatch(callCreateRoomEndpoint());
```

## Redux Slice State

Use slice state for domain state and composed UI state that RTK Query does not already own.

Good slice state:

- current auth user
- auth status
- OAuth callback workflow state
- dashboard composed page model
- dialog open state and dialog payload
- selected filters or view mode

Bad slice state:

- copied mutation loading state
- copied query loading state
- copied API errors used only for one form render
- submitted form values that can stay in React Hook Form or local page state

If RTK Query has the state, do not mirror it in a slice.

## Error Handling

Components and page hooks own user-facing error display.

RTK Query provides the transport error. The component or page hook decides how to show it:

```ts
const message = resolveApiErrorMessage(error, fallbackMessage);
```

Use thunk `rejectWithValue` only when the thunk owns a multi-step workflow and needs to return a domain-level failure for that workflow.

Do not create a thunk only to normalize the error of one endpoint.

## Try/Catch Rules

Components should not use `try/catch` for API calls.

Use RTK Query request state instead.

Rules:

- no `.unwrap()` in components or page hooks
- no bare `try/finally`
- no `try` without `catch`
- no nested `try/catch`

Thunks may use `try/catch` because they own orchestration.

RTK endpoint lifecycle handlers may use `try/catch` when they own side effects such as updating auth state after a successful request.

## Endpoint Lifecycle

Use RTK Query endpoint lifecycle for side effects directly attached to one endpoint.

Examples:

- login/register success sets auth session
- logout success or failure clears local token and auth session
- fetch session success updates auth session
- reset password success clears auth session if required

If lifecycle logic starts coordinating multiple endpoint calls or branching workflow behavior, move that workflow to a thunk.
