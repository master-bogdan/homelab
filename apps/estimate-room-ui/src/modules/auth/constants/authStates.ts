export const AuthRequestStatuses = {
  FAILED: 'failed',
  IDLE: 'idle',
  PENDING: 'pending',
  SUCCEEDED: 'succeeded'
} as const;

export const AuthStates = {
  AUTHENTICATED: 'authenticated',
  UNAUTHENTICATED: 'unauthenticated',
  UNKNOWN: 'unknown'
} as const;
