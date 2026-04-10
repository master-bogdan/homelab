import { Navigate, Outlet, useLocation } from 'react-router-dom';

import { useAppSelector } from '@/app/store/hooks';
import { AUTH_STATUSES } from '@/modules/auth/types';
import { appRoutes } from '@/shared/constants/routes';
import { AppPageState } from '@/shared/ui';
import { selectAuthStatus, selectIsAuthenticated } from '@/modules/auth/store';

export const ProtectedRoute = () => {
  const location = useLocation();
  const isAuthenticated = useAppSelector(selectIsAuthenticated);
  const authStatus = useAppSelector(selectAuthStatus);

  if (authStatus === AUTH_STATUSES.UNKNOWN) {
    return (
      <AppPageState
        description="Checking your current session before opening the workspace."
        isLoading
        minHeight="100vh"
        title="Loading workspace"
        titleComponent="h1"
      />
    );
  }

  if (!isAuthenticated) {
    return <Navigate replace state={{ from: location }} to={appRoutes.login} />;
  }

  return <Outlet />;
};
