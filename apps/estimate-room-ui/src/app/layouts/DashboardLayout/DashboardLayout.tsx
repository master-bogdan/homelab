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

import { DashboardLayoutContent } from './components/DashboardLayoutContent';

export const DashboardLayout = () => {
  const location = useLocation();
  const authStatus = useAppSelector(selectAuthStatus);
  const isAuthenticated = useAppSelector(selectIsAuthenticated);
  const shouldFetchSession = authStatus === AuthStates.UNKNOWN;
  const sessionQuery = useFetchSessionQuery(undefined, {
    refetchOnMountOrArgChange: true,
    skip: !shouldFetchSession
  });
  const isResolvingSession =
    shouldFetchSession &&
    (sessionQuery.isUninitialized || sessionQuery.isLoading || sessionQuery.isFetching);

  if (isResolvingSession) {
    return (
      <AppPageState
        description="Checking your current session before opening the dashboard."
        isLoading
        minHeight="100vh"
        title="Loading session"
        titleComponent="h1"
      />
    );
  }

  if (!isAuthenticated) {
    return <Navigate replace state={{ from: location }} to={AppRoutes.LOGIN} />;
  }

  return <DashboardLayoutContent />;
};
