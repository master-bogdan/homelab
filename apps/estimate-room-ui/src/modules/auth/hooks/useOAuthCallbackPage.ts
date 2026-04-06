import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { useAppDispatch } from '@/app/store/hooks';
import type { AuthUser } from '@/shared/types';

import { authService } from '../services';
import { clearSession, setSession } from '../store';
import type { PendingAuthorizationRequest } from '../types';
import {
  clearPendingAuthorizationRequest,
  readPendingAuthorizationRequest,
  resolveApiErrorMessage
} from '../utils';

type OAuthCallbackResult = {
  readonly pendingRequest: PendingAuthorizationRequest;
  readonly user: AuthUser;
};

const callbackRequests = new Map<string, Promise<OAuthCallbackResult>>();

const createCallbackRequestKey = (code: string, state: string) => `${state}:${code}`;

const finalizeOAuthCallbackRequest = async (
  code: string,
  state: string
): Promise<OAuthCallbackResult> => {
  const pendingRequest = readPendingAuthorizationRequest();

  if (!pendingRequest || pendingRequest.state !== state) {
    throw new Error('Your sign-in session expired. Please try signing in again.');
  }

  await authService.exchangeAuthorizationCode({
    clientId: pendingRequest.clientId,
    code,
    codeVerifier: pendingRequest.codeVerifier,
    redirectUri: pendingRequest.redirectUri
  });

  const user = await authService.fetchSession();

  if (!user) {
    throw new Error('Sign-in completed, but the session could not be restored.');
  }

  return { pendingRequest, user };
};

const getOrCreateCallbackRequest = (code: string, state: string) => {
  const requestKey = createCallbackRequestKey(code, state);
  const inFlightRequest = callbackRequests.get(requestKey);

  if (inFlightRequest) {
    return inFlightRequest;
  }

  const request = finalizeOAuthCallbackRequest(code, state).finally(() => {
    if (callbackRequests.get(requestKey) === request) {
      callbackRequests.delete(requestKey);
    }
  });

  callbackRequests.set(requestKey, request);

  return request;
};

export const useOAuthCallbackPage = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const authorizationCode = searchParams.get('code');
  const authorizationState = searchParams.get('state');
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;

    const finalizeLogin = async () => {
      try {
        const code = authorizationCode;
        const state = authorizationState;

        if (!code || !state) {
          throw new Error('Your sign-in session expired. Please try signing in again.');
        }

        const { pendingRequest, user } = await getOrCreateCallbackRequest(code, state);

        if (!isMounted) {
          return;
        }

        dispatch(setSession(user));
        clearPendingAuthorizationRequest();
        navigate(pendingRequest.redirectTo, { replace: true });
      } catch (error) {
        if (!isMounted) {
          return;
        }

        clearPendingAuthorizationRequest();
        dispatch(clearSession());

        setErrorMessage(
          resolveApiErrorMessage(error, 'Unable to complete sign-in right now.')
        );
      }
    };

    void finalizeLogin();

    return () => {
      isMounted = false;
    };
  }, [authorizationCode, authorizationState, dispatch, navigate]);

  return {
    errorMessage,
    isLoading: errorMessage === null
  };
};
