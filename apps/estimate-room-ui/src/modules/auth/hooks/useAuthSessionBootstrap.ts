import { useEffect } from 'react';

import { useAppDispatch, useAppSelector } from '@/shared/store';
import type { AppDispatch } from '@/shared/store';
import { appRoutes } from '@/shared/constants/routes';
import type { AuthUser } from '@/shared/types';

import { authService } from '../services';
import { selectAuthStatus } from '../store';
import { clearSession, setSession } from '../store';
import { AUTH_STATUSES } from '../types';

let sessionBootstrapRequest: Promise<AuthUser | null> | null = null;

const getOrCreateSessionBootstrapRequest = (dispatch: AppDispatch) => {
  if (sessionBootstrapRequest) {
    return sessionBootstrapRequest;
  }

  const request = (async () => {
    if (!authService.hasStoredAccessToken()) {
      try {
        await authService.refreshAccessToken(dispatch);
      } catch {
        await authService.logout(dispatch).catch(() => undefined);

        return null;
      }
    }

    return authService.fetchSession(dispatch);
  })().finally(() => {
    if (sessionBootstrapRequest === request) {
      sessionBootstrapRequest = null;
    }
  });

  sessionBootstrapRequest = request;

  return request;
};

export const useAuthSessionBootstrap = () => {
  const dispatch = useAppDispatch();
  const authStatus = useAppSelector(selectAuthStatus);

  useEffect(() => {
    if (authStatus !== AUTH_STATUSES.UNKNOWN) {
      return;
    }

    if (
      typeof window !== 'undefined' &&
      window.location.pathname === appRoutes.authCallback
    ) {
      return;
    }

    let isMounted = true;

    const hydrateSession = async () => {
      try {
        const user = await getOrCreateSessionBootstrapRequest(dispatch);

        if (!isMounted) {
          return;
        }

        if (user) {
          dispatch(setSession(user));
          return;
        }

        dispatch(clearSession());
      } catch {
        if (isMounted) {
          dispatch(clearSession());
        }
      }
    };

    void hydrateSession();

    return () => {
      isMounted = false;
    };
  }, [authStatus, dispatch]);
};
