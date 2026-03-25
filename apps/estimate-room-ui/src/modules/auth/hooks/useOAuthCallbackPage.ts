import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { useAppDispatch } from '@/app/store/hooks';

import { authService } from '../services';
import { clearSession, setSession } from '../store';
import {
  clearOauthTokenCookies,
  clearPendingAuthorizationRequest,
  persistOauthTokenCookies,
  readPendingAuthorizationRequest,
  resolveApiErrorMessage
} from '../utils';

export const useOAuthCallbackPage = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;

    const finalizeLogin = async () => {
      try {
        const code = searchParams.get('code');
        const state = searchParams.get('state');
        const pendingRequest = readPendingAuthorizationRequest();

        if (!code || !state || !pendingRequest || pendingRequest.state !== state) {
          throw new Error('Your sign-in session expired. Please try signing in again.');
        }

        const tokens = await authService.exchangeAuthorizationCode({
          clientId: pendingRequest.clientId,
          code,
          codeVerifier: pendingRequest.codeVerifier,
          redirectUri: pendingRequest.redirectUri
        });

        persistOauthTokenCookies(tokens);

        const user = await authService.fetchSession();

        if (!user) {
          throw new Error('Sign-in completed, but the session could not be restored.');
        }

        dispatch(setSession(user));
        clearPendingAuthorizationRequest();
        navigate(pendingRequest.redirectTo, { replace: true });
      } catch (error) {
        clearPendingAuthorizationRequest();
        clearOauthTokenCookies();
        dispatch(clearSession());

        if (isMounted) {
          setErrorMessage(
            resolveApiErrorMessage(error, 'Unable to complete sign-in right now.')
          );
        }
      }
    };

    void finalizeLogin();

    return () => {
      isMounted = false;
    };
  }, [dispatch, navigate, searchParams]);

  return {
    errorMessage,
    isLoading: errorMessage === null
  };
};
