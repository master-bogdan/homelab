import { AppRoutes } from '@/app/router/routePaths';
import { createAppAsyncThunk } from '@/shared/store';
import type { AppDispatch } from '@/shared/types';

import { AuthRequestStatuses } from '../constants';
import type {
  CompleteOAuthCallbackPayload,
  CompleteOAuthCallbackResult,
  OAuthCallbackRequestResult
} from '../types';
import {
  clearPendingAuthorizationRequest,
  readPendingAuthorizationRequest,
  resolveApiErrorMessage
} from '../utils';
import {
  clearSession,
  setOAuthCallbackFailed,
  setOAuthCallbackPending,
  setOAuthCallbackSucceeded,
  setSession
} from './slice';
import { authApi } from '../api/authApi';
import { selectOAuthCallbackState } from './selectors';

const finalizeOAuthCallbackRequest = async (
  dispatch: AppDispatch,
  code: string,
  state: string
): Promise<OAuthCallbackRequestResult> => {
  const pendingRequest = readPendingAuthorizationRequest();

  if (!pendingRequest || pendingRequest.state !== state) {
    throw new Error('Your sign-in session expired. Please try signing in again.');
  }

  const tokenExchangeResult = await dispatch(authApi.endpoints.exchangeAuthorizationCode.initiate({
    clientId: pendingRequest.clientId,
    code,
    codeVerifier: pendingRequest.codeVerifier,
    redirectUri: pendingRequest.redirectUri
  }));

  if (tokenExchangeResult.error) {
    throw tokenExchangeResult.error;
  }

  const sessionResult = await dispatch(authApi.endpoints.fetchSession.initiate(undefined, {
    forceRefetch: true,
    subscribe: false
  }));

  if (sessionResult.error) {
    throw sessionResult.error;
  }

  if (!sessionResult.data) {
    throw new Error('Sign-in completed, but the session could not be restored.');
  }

  return {
    pendingRequest,
    user: sessionResult.data
  };
};

export const completeOAuthCallback = createAppAsyncThunk(
  'auth/completeOAuthCallback',
  async (
    { code, state }: CompleteOAuthCallbackPayload,
    { dispatch, getState, rejectWithValue }
  ): Promise<CompleteOAuthCallbackResult | ReturnType<typeof rejectWithValue>> => {
    const requestKey = `${state}:${code}`;
    const callbackState = selectOAuthCallbackState(getState());

    if (
      callbackState.status === AuthRequestStatuses.SUCCEEDED &&
      callbackState.requestKey === requestKey
    ) {
      return {
        redirectTo: callbackState.redirectTo ?? AppRoutes.DASHBOARD
      };
    }

    if (
      callbackState.status === AuthRequestStatuses.PENDING &&
      callbackState.requestKey === requestKey
    ) {
      return rejectWithValue('Sign-in is already being completed.');
    }

    dispatch(setOAuthCallbackPending(requestKey));

    try {
      const { pendingRequest, user } = await finalizeOAuthCallbackRequest(
        dispatch,
        code,
        state
      );

      dispatch(setSession(user));
      clearPendingAuthorizationRequest();
      dispatch(setOAuthCallbackSucceeded({
        redirectTo: pendingRequest.redirectTo,
        requestKey
      }));

      return {
        redirectTo: pendingRequest.redirectTo
      };
    } catch (error) {
      clearPendingAuthorizationRequest();
      dispatch(clearSession());

      const message = resolveApiErrorMessage(
        error,
        'Unable to complete sign-in right now.'
      );

      dispatch(setOAuthCallbackFailed(message));

      return rejectWithValue(message);
    }
  }
);
