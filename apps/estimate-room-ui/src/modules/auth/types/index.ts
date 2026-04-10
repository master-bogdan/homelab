import type { AuthUser } from '@/shared/types';

export const AUTH_STATUSES = {
  AUTHENTICATED: 'authenticated',
  UNAUTHENTICATED: 'unauthenticated',
  UNKNOWN: 'unknown'
} as const;

export type AuthStatus = (typeof AUTH_STATUSES)[keyof typeof AUTH_STATUSES];
export type ResetPasswordValidationReason = 'expired' | 'invalid' | 'used';

export interface AuthState {
  readonly status: AuthStatus;
  readonly user: AuthUser | null;
}

export interface SessionUserResponseDto {
  readonly avatarUrl?: string | null;
  readonly displayName: string;
  readonly email?: string;
  readonly id: string;
  readonly occupation?: string | null;
  readonly organization?: string | null;
}

export interface SessionResponseDto {
  readonly authenticated: boolean;
  readonly user: SessionUserResponseDto | null;
}

export interface LoginPayload {
  readonly continue: string;
  readonly email: string;
  readonly password: string;
}

export interface RegisterPayload {
  readonly continue: string;
  readonly displayName: string;
  readonly email: string;
  readonly occupation?: string;
  readonly organization?: string;
  readonly password: string;
}

export interface ForgotPasswordPayload {
  readonly email: string;
}

export interface ResetPasswordPayload {
  readonly password: string;
  readonly token: string;
}

export interface ResetPasswordValidationResponseDto {
  readonly reason?: ResetPasswordValidationReason;
  readonly valid: boolean;
}

export interface OAuthTokenResponseDto {
  readonly access_token: string;
  readonly expires_in: number;
  readonly id_token: string;
  readonly refresh_token?: string;
  readonly token_type: string;
}

export interface OAuthTokenResponse {
  readonly accessToken: string;
  readonly expiresIn: number;
  readonly idToken: string;
  readonly refreshToken?: string;
  readonly tokenType: string;
}

export interface PendingAuthorizationRequest {
  readonly clientId: string;
  readonly codeVerifier: string;
  readonly continueUrl: string;
  readonly redirectTo: string;
  readonly redirectUri: string;
  readonly state: string;
}
