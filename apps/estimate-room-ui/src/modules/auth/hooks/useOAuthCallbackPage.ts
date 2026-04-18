import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { useAppDispatch, useAppSelector } from '@/shared/hooks';

import { AuthRequestStatuses } from '../constants';
import { clearSession, completeOAuthCallback, selectOAuthCallbackState } from '../store';
import { resolveApiErrorMessage } from '../utils';

const getOAuthCallbackErrorMessage = (
  isMissingAuthorizationParams: boolean,
  errorMessage: string | null,
  status: string
) => {
  if (isMissingAuthorizationParams) {
    return 'Your sign-in session expired. Please try signing in again.';
  }

  if (status !== AuthRequestStatuses.FAILED) {
    return null;
  }

  return resolveApiErrorMessage(
    errorMessage,
    'Unable to complete sign-in right now.'
  );
};

export const useOAuthCallbackPage = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const authorizationCode = searchParams.get('code');
  const authorizationState = searchParams.get('state');
  const oauthCallback = useAppSelector(selectOAuthCallbackState);
  const isMissingAuthorizationParams = !authorizationCode || !authorizationState;
  const errorMessage = getOAuthCallbackErrorMessage(
    isMissingAuthorizationParams,
    oauthCallback.errorMessage,
    oauthCallback.status
  );

  useEffect(() => {
    if (isMissingAuthorizationParams) {
      dispatch(clearSession());
      return;
    }

    if (oauthCallback.status !== AuthRequestStatuses.IDLE) {
      return;
    }

    dispatch(completeOAuthCallback({
      code: authorizationCode,
      state: authorizationState
    }));
  }, [
    authorizationCode,
    authorizationState,
    dispatch,
    isMissingAuthorizationParams,
    oauthCallback.status
  ]);

  useEffect(() => {
    if (
      oauthCallback.status !== AuthRequestStatuses.SUCCEEDED ||
      !oauthCallback.redirectTo
    ) {
      return;
    }

    navigate(oauthCallback.redirectTo, { replace: true });
  }, [navigate, oauthCallback.redirectTo, oauthCallback.status]);

  return {
    errorMessage,
    isLoading: errorMessage === null
  };
};
