import { api, accessTokenStorage } from '@/shared/api';
import type { AuthUser } from '@/shared/types';

import type {
  ForgotPasswordPayload,
  LoginPayload,
  OAuthTokenApiResponse,
  OAuthTokenResponse,
  RegisterPayload,
  ResetPasswordValidationApiResponse,
  ResetPasswordPayload,
  SessionApiResponse,
  SessionUserApiResponse
} from '../types';
import { clearSession, setSession } from './authSlice';

const mapSessionUser = (user: SessionUserApiResponse | null): AuthUser | null => {
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

const mapRequiredUser = (response: SessionApiResponse) => {
  const user = response.authenticated ? mapSessionUser(response.user) : null;

  if (!user) {
    throw new Error('Authentication did not return an active session.');
  }

  return user;
};

const mapTokenResponse = (response: OAuthTokenApiResponse): OAuthTokenResponse => ({
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
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;

          if (data) {
            dispatch(setSession(data));
          } else {
            dispatch(clearSession());
          }
        } catch {
          dispatch(clearSession());
        }
      },
      query: () => ({
        url: 'auth/session'
      }),
      transformResponse: (response: SessionApiResponse) =>
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
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;

          dispatch(setSession(data));
        } catch {
          return;
        }
      },
      query: (payload) => ({
        body: payload,
        method: 'POST',
        url: 'auth/login'
      }),
      transformResponse: mapRequiredUser
    }),
    logout: builder.mutation<{ loggedOut: boolean }, void>({
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          await queryFulfilled;
        } catch {
          // Local cleanup still runs when the server logout request fails.
        } finally {
          accessTokenStorage.clear();
          dispatch(api.util.resetApiState());
          dispatch(clearSession());
        }
      },
      query: () => ({
        method: 'POST',
        url: 'auth/logout'
      })
    }),
    register: builder.mutation<AuthUser, RegisterPayload>({
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;

          dispatch(setSession(data));
        } catch {
          return;
        }
      },
      query: (payload) => ({
        body: payload,
        method: 'POST',
        url: 'auth/register'
      }),
      transformResponse: mapRequiredUser
    }),
    resetPassword: builder.mutation<{ reset: boolean }, ResetPasswordPayload>({
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          await queryFulfilled;

          dispatch(clearSession());
        } catch {
          return;
        }
      },
      query: (payload) => ({
        body: payload,
        method: 'POST',
        url: 'auth/reset-password'
      })
    }),
    validateResetPasswordToken: builder.query<ResetPasswordValidationApiResponse, string>({
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
  useFetchSessionQuery,
  useForgotPasswordMutation,
  useLazyValidateResetPasswordTokenQuery,
  useLoginMutation,
  useLogoutMutation,
  useRegisterMutation,
  useResetPasswordMutation,
  useValidateResetPasswordTokenQuery
} = authApi;
