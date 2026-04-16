import type { RouteObject } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';

import { DashboardPage } from './DashboardPage';
import { JoinRoomPage } from './JoinRoomPage';

export const dashboardRoutes: RouteObject[] = [
  {
    path: AppRoutes.DASHBOARD,
    element: <DashboardPage />
  },
  {
    path: AppRoutes.JOIN_ROOM,
    element: <JoinRoomPage />
  }
];
