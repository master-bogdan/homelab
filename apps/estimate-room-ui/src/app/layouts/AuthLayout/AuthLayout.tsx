import { Navigate, useLocation } from 'react-router-dom';

import { useAppSelector } from '@/shared/store';
import { AUTH_STATUSES, selectAuthStatus, selectIsAuthenticated } from '@/modules/auth';
import { appRoutes } from '@/shared/constants/routes';
import { AppPageState } from '@/shared/ui';

import { AuthLayoutContent } from './AuthLayoutContent';
import { resolveAuthRedirectTarget, type RedirectStateLike } from './authLayout.utils';

export const AuthLayout = () => {
  const location = useLocation();
  const authStatus = useAppSelector(selectAuthStatus);
  const isAuthenticated = useAppSelector(selectIsAuthenticated);

  if (location.pathname === appRoutes.authCallback) {
    return <AuthLayoutContent />;
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
        to={resolveAuthRedirectTarget(location.state as RedirectStateLike | null)}
      />
    );
  }

  return <AuthLayoutContent />;
};
