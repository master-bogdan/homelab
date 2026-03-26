import { Navigate, Outlet, useLocation } from 'react-router-dom';

import { useAppSelector } from '@/app/store/hooks';
import { appRoutes } from '@/shared/constants/routes';
import { AppPageState } from '@/shared/ui';
import { selectAuthStatus, selectIsAuthenticated } from '@/modules/auth/selectors';

export const ProtectedRoute = () => {
  const location = useLocation();
  const isAuthenticated = useAppSelector(selectIsAuthenticated);
  const authStatus = useAppSelector(selectAuthStatus);

  if (authStatus === 'unknown') {
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
