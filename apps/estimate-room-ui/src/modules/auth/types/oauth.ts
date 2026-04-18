import type { AuthUser } from '@/modules/auth/types';

export interface CompleteOAuthCallbackPayload {
  readonly code: string;
  readonly state: string;
}

export interface CompleteOAuthCallbackResult {
  readonly redirectTo: string;
}

export interface PendingAuthorizationRequest {
  readonly clientId: string;
  readonly codeVerifier: string;
  readonly continueUrl: string;
  readonly redirectTo: string;
  readonly redirectUri: string;
  readonly state: string;
}

export interface OAuthCallbackRequestResult {
  readonly pendingRequest: PendingAuthorizationRequest;
  readonly user: AuthUser;
}
