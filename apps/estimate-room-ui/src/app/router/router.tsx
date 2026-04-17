import { Navigate, createBrowserRouter } from 'react-router-dom';

import { AuthLayout, DashboardLayout } from '@/app/layouts';
import { AuthRoutes } from '@/modules/auth';
import { DashboardRoutes } from '@/modules/dashboard';
import { HistoryRoutes } from '@/modules/history';
import { ProfileRoutes } from '@/modules/profile';
import { RoomsRoutes } from '@/modules/rooms';
import { SettingsRoutes } from '@/modules/settings';
import { TeamsRoutes } from '@/modules/teams';
import { AppRoutes } from '@/shared/constants/routes';

import { NotFoundPage } from './NotFoundPage';

export const router = createBrowserRouter([
  {
    path: AppRoutes.ROOT,
    element: <Navigate replace to={AppRoutes.DASHBOARD} />
  },
  {
    element: <AuthLayout />,
    children: AuthRoutes
  },
  {
    element: <DashboardLayout />,
    children: [
      ...DashboardRoutes,
      ...RoomsRoutes,
      ...HistoryRoutes,
      ...TeamsRoutes,
      ...ProfileRoutes,
      ...SettingsRoutes
    ]
  },
  {
    path: '*',
    element: <NotFoundPage />
  }
]);
