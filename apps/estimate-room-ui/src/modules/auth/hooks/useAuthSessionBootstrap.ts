import { useEffect } from 'react';

import { useAppDispatch, useAppSelector } from '@/shared/store';
import { appRoutes } from '@/shared/constants/routes';

import { bootstrapAuthSession, clearSession, selectAuthStatus } from '../store';
import { AUTH_STATUSES } from '../types';

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
        await dispatch(bootstrapAuthSession()).unwrap();
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
