import { AppConfig } from '@/config';
import { api, accessTokenStorage } from '@/shared/api';
import type { AuthUser } from '@/shared/types';

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

export const authApi = api.injectEndpoints({
  endpoints: (builder) => ({
    exchangeAuthorizationCode: builder.mutation<
      OAuthTokenResponse,
      {
        readonly clientId: string;
        readonly code: string;
        readonly codeVerifier: string;
        readonly redirectUri: string;
      }
    >({
      async onQueryStarted(_arg, { queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;

          accessTokenStorage.set(data.accessToken);
        } catch {
          return;
        }
      },
      query: ({ clientId, code, codeVerifier, redirectUri }) => ({
        body: new URLSearchParams({
          client_id: clientId,
          code,
          code_verifier: codeVerifier,
          grant_type: 'authorization_code',
          redirect_uri: redirectUri
        }).toString(),
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        },
        method: 'POST',
        url: 'oauth2/token'
      }),
      transformResponse: mapTokenResponse
    }),
    fetchSession: builder.query<AuthUser | null, void>({
      query: () => ({
        url: 'auth/session'
      }),
      transformResponse: (response: SessionResponseDto) =>
        response.authenticated ? mapSessionUser(response.user) : null
    }),
    forgotPassword: builder.mutation<{ submitted: boolean }, ForgotPasswordPayload>({
      query: (payload) => ({
        body: payload,
        method: 'POST',
        url: 'auth/forgot-password'
      })
    }),
    login: builder.mutation<AuthUser, LoginPayload>({
      query: (payload) => ({
        body: payload,
        method: 'POST',
        url: 'auth/login'
      }),
      transformResponse: mapRequiredUser
    }),
    logout: builder.mutation<{ loggedOut: boolean }, void>({
      async onQueryStarted(_arg, { queryFulfilled }) {
        try {
          await queryFulfilled;
        } finally {
          accessTokenStorage.clear();
        }
      },
      query: () => ({
        method: 'POST',
        url: 'auth/logout'
      })
    }),
    refreshAccessToken: builder.mutation<OAuthTokenResponse, void>({
      async onQueryStarted(_arg, { queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;

          accessTokenStorage.set(data.accessToken);
        } catch {
          return;
        }
      },
      query: () => ({
        body: new URLSearchParams({
          client_id: AppConfig.OAUTH_CLIENT_ID.trim(),
          grant_type: 'refresh_token'
        }).toString(),
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        },
        method: 'POST',
        url: 'oauth2/token'
      }),
      transformResponse: mapTokenResponse
    }),
    register: builder.mutation<AuthUser, RegisterPayload>({
      query: (payload) => ({
        body: payload,
        method: 'POST',
        url: 'auth/register'
      }),
      transformResponse: mapRequiredUser
    }),
    resetPassword: builder.mutation<{ reset: boolean }, ResetPasswordPayload>({
      query: (payload) => ({
        body: payload,
        method: 'POST',
        url: 'auth/reset-password'
      })
    }),
    validateResetPasswordToken: builder.query<ResetPasswordValidationResponseDto, string>({
      query: (token) => ({
        params: {
          token
        },
        url: 'auth/reset-password/validate'
      })
    })
  }),
  overrideExisting: false
});

export const {
  useForgotPasswordMutation,
  useLazyValidateResetPasswordTokenQuery,
  useResetPasswordMutation,
  useValidateResetPasswordTokenQuery
} = authApi;
