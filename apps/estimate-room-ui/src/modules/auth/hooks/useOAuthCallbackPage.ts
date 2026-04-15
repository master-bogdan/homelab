import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { useAppDispatch } from '@/shared/store';

import { clearSession, completeOAuthCallback } from '../store';
import { resolveApiErrorMessage } from '../utils';

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

        const { redirectTo } = await dispatch(
          completeOAuthCallback({ code, state })
        ).unwrap();

        if (!isMounted) {
          return;
        }

        navigate(redirectTo, { replace: true });
      } catch (error) {
        if (!isMounted) {
          return;
        }

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
