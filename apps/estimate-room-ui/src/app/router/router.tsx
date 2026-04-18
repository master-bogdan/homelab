import { Navigate, createBrowserRouter } from 'react-router-dom';

import { AuthLayout, DashboardLayout } from '@/app/layouts';
import { NotFoundPage } from '@/app/pages';
import { AppRoutes } from '@/app/router/routePaths';

import {
  AuthRoutes,
  DashboardRoutes,
  HistoryRoutes,
  ProfileRoutes,
  RoomsRoutes,
  SettingsRoutes,
  TeamsRoutes
} from './routes';

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
