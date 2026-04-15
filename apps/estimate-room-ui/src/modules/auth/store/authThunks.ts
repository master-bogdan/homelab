import { api, accessTokenStorage } from '@/shared/api';
import { createAppAsyncThunk } from '@/shared/store';
import type { AppDispatch } from '@/shared/store';
import type { AuthUser } from '@/shared/types';

import type {
  LoginPayload,
  PendingAuthorizationRequest,
  RegisterPayload,
  ResetPasswordPayload
} from '../types';
import {
  clearPendingAuthorizationRequest,
  isEmailAlreadyInUseError,
  isInvalidCredentialsError,
  readPendingAuthorizationRequest,
  resolveApiErrorMessage
} from '../utils';
import { clearSession, setSession } from './authSlice';
import { authApi } from './authService';

type CompleteOAuthCallbackPayload = {
  readonly code: string;
  readonly state: string;
};

type CompleteOAuthCallbackResult = {
  readonly redirectTo: string;
};

type OAuthCallbackRequestResult = {
  readonly pendingRequest: PendingAuthorizationRequest;
  readonly user: AuthUser;
};

let sessionBootstrapRequest: Promise<AuthUser | null> | null = null;
const callbackRequests = new Map<string, Promise<OAuthCallbackRequestResult>>();

const createCallbackRequestKey = (code: string, state: string) => `${state}:${code}`;

export const submitLogin = createAppAsyncThunk(
  'auth/submitLogin',
  async (payload: LoginPayload, { dispatch, rejectWithValue }) => {
    try {
      const user = await dispatch(authApi.endpoints.login.initiate(payload)).unwrap();

      dispatch(setSession(user));

      return user;
    } catch (error) {
      return rejectWithValue(
        isInvalidCredentialsError(error)
          ? 'Email or password is incorrect.'
          : resolveApiErrorMessage(error, 'Unable to sign in right now.')
      );
    }
  }
);

export const submitRegister = createAppAsyncThunk(
  'auth/submitRegister',
  async (payload: RegisterPayload, { dispatch, rejectWithValue }) => {
    try {
      const user = await dispatch(authApi.endpoints.register.initiate(payload)).unwrap();

      dispatch(setSession(user));

      return user;
    } catch (error) {
      return rejectWithValue(
        isEmailAlreadyInUseError(error)
          ? 'This email is already registered.'
          : resolveApiErrorMessage(error, 'Unable to create your account right now.')
      );
    }
  }
);

export const bootstrapAuthSession = createAppAsyncThunk(
  'auth/bootstrapAuthSession',
  async (_, { dispatch }) => {
    if (!sessionBootstrapRequest) {
      const request = (async () => {
        if (!accessTokenStorage.get()) {
          try {
            await dispatch(authApi.endpoints.refreshAccessToken.initiate()).unwrap();
          } catch {
            await dispatch(authApi.endpoints.logout.initiate()).unwrap().catch(() => undefined);

            return null;
          }
        }

        return dispatch(authApi.endpoints.fetchSession.initiate(undefined, {
          forceRefetch: true,
          subscribe: false
        })).unwrap();
      })().finally(() => {
        if (sessionBootstrapRequest === request) {
          sessionBootstrapRequest = null;
        }
      });

      sessionBootstrapRequest = request;
    }

    const user = await sessionBootstrapRequest;

    if (user) {
      dispatch(setSession(user));
    } else {
      dispatch(clearSession());
    }

    return user;
  }
);

const finalizeOAuthCallbackRequest = async (
  dispatch: AppDispatch,
  code: string,
  state: string
): Promise<OAuthCallbackRequestResult> => {
  const pendingRequest = readPendingAuthorizationRequest();

  if (!pendingRequest || pendingRequest.state !== state) {
    throw new Error('Your sign-in session expired. Please try signing in again.');
  }

  await dispatch(authApi.endpoints.exchangeAuthorizationCode.initiate({
    clientId: pendingRequest.clientId,
    code,
    codeVerifier: pendingRequest.codeVerifier,
    redirectUri: pendingRequest.redirectUri
  })).unwrap();

  const user = await dispatch(authApi.endpoints.fetchSession.initiate(undefined, {
    forceRefetch: true,
    subscribe: false
  })).unwrap();

  if (!user) {
    throw new Error('Sign-in completed, but the session could not be restored.');
  }

  return { pendingRequest, user };
};

export const completeOAuthCallback = createAppAsyncThunk(
  'auth/completeOAuthCallback',
  async (
    { code, state }: CompleteOAuthCallbackPayload,
    { dispatch, rejectWithValue }
  ): Promise<CompleteOAuthCallbackResult | ReturnType<typeof rejectWithValue>> => {
    try {
      const requestKey = createCallbackRequestKey(code, state);
      const inFlightRequest = callbackRequests.get(requestKey);

      if (!inFlightRequest) {
        const request = finalizeOAuthCallbackRequest(dispatch, code, state).finally(() => {
          if (callbackRequests.get(requestKey) === request) {
            callbackRequests.delete(requestKey);
          }
        });

        callbackRequests.set(requestKey, request);
      }

      const { pendingRequest, user } = await callbackRequests.get(requestKey)!;

      dispatch(setSession(user));
      clearPendingAuthorizationRequest();

      return {
        redirectTo: pendingRequest.redirectTo
      };
    } catch (error) {
      clearPendingAuthorizationRequest();
      dispatch(clearSession());

      return rejectWithValue(
        resolveApiErrorMessage(error, 'Unable to complete sign-in right now.')
      );
    }
  }
);

export const submitLogout = createAppAsyncThunk(
  'auth/submitLogout',
  async (_, { dispatch }) => {
    try {
      await dispatch(authApi.endpoints.logout.initiate()).unwrap();
    } finally {
      dispatch(api.util.resetApiState());
      dispatch(clearSession());
    }

    return { loggedOut: true };
  }
);

export const submitResetPassword = createAppAsyncThunk(
  'auth/submitResetPassword',
  async (payload: ResetPasswordPayload, { dispatch, rejectWithValue }) => {
    try {
      const result = await dispatch(authApi.endpoints.resetPassword.initiate(payload)).unwrap();

      dispatch(clearSession());

      return result;
    } catch (error) {
      return rejectWithValue(
        resolveApiErrorMessage(error, 'Unable to reset your password right now.')
      );
    }
  }
);
