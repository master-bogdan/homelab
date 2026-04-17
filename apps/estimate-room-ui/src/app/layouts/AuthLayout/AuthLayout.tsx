import { Navigate, useLocation } from 'react-router-dom';

import { useAppSelector } from '@/shared/store';
import {
  AuthStates,
  selectAuthStatus,
  selectIsAuthenticated,
  useFetchSessionQuery
} from '@/modules/auth';
import { AppRoutes } from '@/shared/constants/routes';
import { AppPageState } from '@/shared/ui';

import { AuthLayoutContent } from './components';
import { resolveAuthRedirectTarget, type RedirectStateLike } from './utils';

export const AuthLayout = () => {
  const location = useLocation();
  const authStatus = useAppSelector(selectAuthStatus);
  const isAuthenticated = useAppSelector(selectIsAuthenticated);
  const shouldFetchSession =
    authStatus === AuthStates.UNKNOWN && location.pathname !== AppRoutes.AUTH_CALLBACK;
  const sessionQuery = useFetchSessionQuery(undefined, {
    refetchOnMountOrArgChange: true,
    skip: !shouldFetchSession
  });
  const isResolvingSession =
    shouldFetchSession &&
    (sessionQuery.isUninitialized || sessionQuery.isLoading || sessionQuery.isFetching);

  if (location.pathname === AppRoutes.AUTH_CALLBACK) {
    return <AuthLayoutContent />;
  }

  if (isResolvingSession) {
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
