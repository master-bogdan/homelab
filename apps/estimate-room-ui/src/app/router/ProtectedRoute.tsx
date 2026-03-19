import { Box, CircularProgress } from '@mui/material';
import { Navigate, Outlet, useLocation } from 'react-router-dom';

import { useAppSelector } from '@/app/store/hooks';
import { appRoutes } from '@/shared/constants/routes';
import { selectAuthStatus, selectIsAuthenticated } from '@/modules/auth/selectors';

export const ProtectedRoute = () => {
  const location = useLocation();
  const isAuthenticated = useAppSelector(selectIsAuthenticated);
  const authStatus = useAppSelector(selectAuthStatus);

  if (authStatus === 'unknown') {
    return (
      <Box
        sx={{
          display: 'grid',
          placeItems: 'center',
          minHeight: '100vh'
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  if (!isAuthenticated) {
    return <Navigate replace state={{ from: location.pathname }} to={appRoutes.login} />;
  }

  return <Outlet />;
};
