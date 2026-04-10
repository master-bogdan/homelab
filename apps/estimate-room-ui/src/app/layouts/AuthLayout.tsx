import { Navigate, Outlet, useLocation } from 'react-router-dom';

import { useAppSelector } from '@/app/store/hooks';
import { selectAuthStatus, selectIsAuthenticated } from '@/modules/auth/store';
import { AUTH_STATUSES } from '@/modules/auth/types';
import { appRoutes } from '@/shared/constants/routes';
import { AppPageState } from '@/shared/ui';

interface RedirectStateLike {
  readonly from?: {
    readonly hash: string;
    readonly pathname: string;
    readonly search: string;
  };
}

const resolveRedirectTarget = (state: RedirectStateLike | null) => {
  const from = state?.from;

  if (!from) {
    return appRoutes.dashboard;
  }

  return `${from.pathname}${from.search}${from.hash}`;
};

export const AuthLayout = () => {
  const location = useLocation();
  const authStatus = useAppSelector(selectAuthStatus);
  const isAuthenticated = useAppSelector(selectIsAuthenticated);

  if (location.pathname === appRoutes.authCallback) {
    return <Outlet />;
  }

  if (authStatus === AUTH_STATUSES.UNKNOWN) {
    return (
      <AppPageState
        description="Checking your current session before opening authentication pages."
        isLoading
        minHeight="100vh"
        title="Loading session"
        titleComponent="h1"
      />
    );
  }

  if (isAuthenticated) {
    return (
      <Navigate
        replace
        to={resolveRedirectTarget(location.state as RedirectStateLike | null)}
      />
    );
  }

  return <Outlet />;
};
