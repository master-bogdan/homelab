import { useEffect } from 'react';

import { useAppDispatch, useAppSelector } from '@/app/store/hooks';

import { authService } from '../services';
import { selectAuthStatus } from '../selectors';
import { clearSession, setSession } from '../store';

export const useAuthSessionBootstrap = () => {
  const dispatch = useAppDispatch();
  const authStatus = useAppSelector(selectAuthStatus);

  useEffect(() => {
    if (authStatus !== 'unknown') {
      return;
    }

    let isMounted = true;

    const hydrateSession = async () => {
      try {
        const user = await authService.fetchSession();

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
