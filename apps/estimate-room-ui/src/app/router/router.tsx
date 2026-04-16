import { Navigate, createBrowserRouter } from 'react-router-dom';

import { AuthLayout, DashboardLayout } from '@/app/layouts';
import { authRoutes } from '@/modules/auth';
import { dashboardRoutes } from '@/modules/dashboard';
import { historyRoutes } from '@/modules/history';
import { profileRoutes } from '@/modules/profile';
import { roomsRoutes } from '@/modules/rooms';
import { settingsRoutes } from '@/modules/settings';
import { teamsRoutes } from '@/modules/teams';
import { AppRoutes } from '@/shared/constants/routes';

import { NotFoundPage } from './NotFoundPage';

export const router = createBrowserRouter([
  {
    path: AppRoutes.ROOT,
    element: <Navigate replace to={AppRoutes.DASHBOARD} />
  },
  {
    element: <AuthLayout />,
    children: authRoutes
  },
  {
    element: <DashboardLayout />,
    children: [
      ...dashboardRoutes,
      ...roomsRoutes,
      ...historyRoutes,
      ...teamsRoutes,
      ...profileRoutes,
      ...settingsRoutes
    ]
  },
  {
    path: '*',
    element: <NotFoundPage />
  }
]);
