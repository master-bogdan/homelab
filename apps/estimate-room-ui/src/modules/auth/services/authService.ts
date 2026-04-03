import { apiClient } from '@/shared/api';
import type { AuthUser } from '@/shared/types';

import { createApiUrl } from '../utils';
import type {
  ForgotPasswordPayload,
  LoginPayload,
  OAuthTokenResponse,
  OAuthTokenResponseDto,
  RegisterPayload,
  ResetPasswordPayload,
  ResetPasswordValidationResponseDto,
  SessionResponseDto,
  SessionUserResponseDto
} from '../types';

const mapSessionUser = (user: SessionUserResponseDto | null): AuthUser | null => {
  if (!user) {
    return null;
  }

  return {
    avatarUrl: user.avatarUrl ?? null,
    displayName: user.displayName,
    email: user.email ?? '',
    id: user.id,
    occupation: user.occupation ?? null,
    organization: user.organization ?? null
  };
};

const mapRequiredUser = (response: SessionResponseDto) => {
  const user = response.authenticated ? mapSessionUser(response.user) : null;

  if (!user) {
    throw new Error('Authentication did not return an active session.');
  }

  return user;
};

const mapTokenResponse = (response: OAuthTokenResponseDto): OAuthTokenResponse => ({
  accessToken: response.access_token,
  expiresIn: response.expires_in,
  idToken: response.id_token,
  refreshToken: response.refresh_token,
  tokenType: response.token_type
});

export const authService = {
  exchangeAuthorizationCode: async ({
    clientId,
    code,
    codeVerifier,
    redirectUri
  }: {
    readonly clientId: string;
    readonly code: string;
    readonly codeVerifier: string;
    readonly redirectUri: string;
  }) => {
    const form = new URLSearchParams({
      client_id: clientId,
      code,
      code_verifier: codeVerifier,
      grant_type: 'authorization_code',
      redirect_uri: redirectUri
    });

    const response = await apiClient.post<OAuthTokenResponseDto>('oauth2/token', form, {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    });

    return mapTokenResponse(response);
  },
  fetchSession: async () => {
    const response = await apiClient.get<SessionResponseDto>('auth/session');

    return response.authenticated ? mapSessionUser(response.user) : null;
  },
  forgotPassword: async (payload: ForgotPasswordPayload) =>
    apiClient.post<{ submitted: boolean }>('auth/forgot-password', payload),
  getGithubLoginUrl: (continueUrl: string) =>
    createApiUrl('auth/github/login', { continue: continueUrl }).toString(),
  login: async (payload: LoginPayload) => {
    const response = await apiClient.post<SessionResponseDto>('auth/login', payload);

    return mapRequiredUser(response);
  },
  logout: async () => apiClient.post<{ loggedOut: boolean }>('auth/logout'),
  register: async (payload: RegisterPayload) => {
    const response = await apiClient.post<SessionResponseDto>('auth/register', payload);

    return mapRequiredUser(response);
  },
  resetPassword: async (payload: ResetPasswordPayload) =>
    apiClient.post<{ reset: boolean }>('auth/reset-password', payload),
  validateResetPasswordToken: async (token: string) =>
    apiClient.get<ResetPasswordValidationResponseDto>('auth/reset-password/validate', {
      query: { token }
    })
};
