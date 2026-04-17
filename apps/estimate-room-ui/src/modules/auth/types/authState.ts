import type { AuthUser } from '@/shared/types';

import type { AuthStates } from '../constants';

export type AuthStatus = (typeof AuthStates)[keyof typeof AuthStates];
export type AuthRequestStatus = 'failed' | 'idle' | 'pending' | 'succeeded';

export interface OAuthCallbackState {
  readonly errorMessage: string | null;
  readonly redirectTo: string | null;
  readonly requestKey: string | null;
  readonly status: AuthRequestStatus;
}

export interface AuthState {
  readonly oauthCallback: OAuthCallbackState;
  readonly status: AuthStatus;
  readonly user: AuthUser | null;
}
