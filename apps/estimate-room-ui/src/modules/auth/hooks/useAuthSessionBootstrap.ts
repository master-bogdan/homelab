import { useEffect } from 'react';

import { useAppDispatch, useAppSelector } from '@/app/store/hooks';
import { appRoutes } from '@/shared/constants/routes';
import type { AuthUser } from '@/shared/types';

import { authService } from '../services';
import { selectAuthStatus } from '../selectors';
import { clearSession, setSession } from '../store';

let sessionBootstrapRequest: Promise<AuthUser | null> | null = null;

const getOrCreateSessionBootstrapRequest = () => {
  if (sessionBootstrapRequest) {
    return sessionBootstrapRequest;
  }

  const request = authService.fetchSession().finally(() => {
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
    if (authStatus !== 'unknown') {
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
        const user = await getOrCreateSessionBootstrapRequest();

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
