import type { ResetPasswordValidationReason } from './resetPassword';

export interface OAuthTokenApiResponse {
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

export interface ResetPasswordValidationApiResponse {
  readonly reason?: ResetPasswordValidationReason;
  readonly valid: boolean;
}

export interface SessionUserApiResponse {
  readonly avatarUrl?: string | null;
  readonly displayName: string;
  readonly email?: string;
  readonly id: string;
  readonly occupation?: string | null;
  readonly organization?: string | null;
}

export interface SessionApiResponse {
  readonly authenticated: boolean;
  readonly user: SessionUserApiResponse | null;
}
