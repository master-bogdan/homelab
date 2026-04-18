# Request Flow Decisions

This table is the current decision map for auth and dashboard requests.

## Auth

| Flow | Owner | Reason |
| --- | --- | --- |
| Fetch current session | RTK Query hook | One endpoint. Layout uses query state while auth is unknown. |
| Login | RTK Query mutation hook | One endpoint. Endpoint lifecycle can set session. Component owns form error display. |
| Register | RTK Query mutation hook | One endpoint. Endpoint lifecycle can set session. Component owns form error display. |
| Logout | RTK Query mutation hook | One endpoint. Endpoint lifecycle can clear token, reset API cache, and clear session. |
| Forgot password | RTK Query mutation hook inside focused flow hook | One endpoint. `useForgotPasswordPage` owns submitted email state, resend behavior, and error normalization. |
| Resend forgot password | Same focused flow hook | Same endpoint as forgot password. No separate thunk. |
| Validate reset token | RTK Query query hook | One endpoint. Page renders query state. |
| Reset password | RTK Query mutation hook | One endpoint. Endpoint lifecycle can clear session if needed. |
| OAuth callback | Thunk | Multi-step workflow: read pending request, exchange code, fetch session, update state, clear pending request. |

## Dashboard

| Flow | Owner | Reason |
| --- | --- | --- |
| Dashboard page load | Thunk | Composes sessions, teams, ledger, active room, errors, and view state. |
| Create room | Thunk backed by RTK Query mutation | Creates a room, refreshes dashboard page state, and lets UI open the success dialog from the thunk result. |
| Refresh dashboard after create room | Dashboard thunk or RTK invalidation | Refresh is dashboard page state, not create-room request ownership. |
| Join room | Thunk | Multi-step workflow: parse invite, preview invite, reject wrong kind, accept room invite. |
| Fetch create-room teams | Thunk backed by RTK Query query | Loads team options into dashboard dialog state for the existing dialog workflow. |
| Preview invitation only | RTK Query query hook | One endpoint if used alone. |

## Decision Checklist

Before adding a request, answer these questions:

1. Is this exactly one API endpoint?
   - Yes: use RTK Query hook.
   - No: use a thunk.

2. Does the UI only need loading, data, and error?
   - Yes: use RTK Query state directly.

3. Does the action need multiple endpoint calls or conditional branching?
   - Yes: use a thunk.

4. Does the result become long-lived domain state?
   - Yes: store the domain state in a slice.
   - No: keep it in RTK Query, React Hook Form, or local component state.

5. Are you creating slice state that duplicates RTK Query request state?
   - If yes, stop and use RTK Query state.

## Naming

Use domain names for thunks:

- `fetchDashboardPage`
- `completeOAuthCallback`
- `submitJoinRoom`

Avoid transport names for thunks:

- `submitLogout`
- `callCreateRoomEndpoint`
- `sendForgotPassword`

Those belong to RTK Query mutation hooks unless they coordinate a larger workflow.

## Page Hook Decision

Use page hooks only for real page orchestration.

Current auth decisions:

- `useLoginPage`: keep. It combines form setup, auth continuation, GitHub redirect,
  RTK login result handling, and final redirect.
- `useRegisterPage`: keep. It combines form setup, password confirmation
  revalidation, auth continuation, GitHub redirect, RTK register result handling,
  and final redirect.
- `useResetPasswordPage`: keep. It combines token query-param handling, token
  validation, page state, reset form state, RTK mutation handling, and navigation.
- `useOAuthCallbackPage`: keep. It connects router state to the
  `completeOAuthCallback` thunk and navigation.
- `useForgotPasswordPage`: keep. It combines form setup, mutation submission,
  submitted email state, resend behavior, and request error normalization.
