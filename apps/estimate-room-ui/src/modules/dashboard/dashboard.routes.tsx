import type { RouteObject } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';

import { DashboardPage } from './DashboardPage';
import { JoinRoomPage } from './JoinRoomPage';

export const dashboardRoutes: RouteObject[] = [
  {
    path: appRoutes.dashboard,
    element: <DashboardPage />
  },
  {
    path: appRoutes.joinRoom,
    element: <JoinRoomPage />
  }
];
