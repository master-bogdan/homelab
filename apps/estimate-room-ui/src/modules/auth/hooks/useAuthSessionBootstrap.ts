import { useEffect } from 'react';

import { useAppDispatch, useAppSelector } from '@/shared/store';
import { AppRoutes } from '@/shared/constants/routes';

import { bootstrapAuthSession, clearSession, selectAuthStatus } from '../store';
import { AuthStates } from '../types';

export const useAuthSessionBootstrap = () => {
  const dispatch = useAppDispatch();
  const authStatus = useAppSelector(selectAuthStatus);

  useEffect(() => {
    if (authStatus !== AuthStates.UNKNOWN) {
      return;
    }

    if (
      typeof window !== 'undefined' &&
      window.location.pathname === AppRoutes.AUTH_CALLBACK
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
