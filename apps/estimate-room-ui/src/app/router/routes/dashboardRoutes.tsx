import type { RouteObject } from 'react-router-dom';

import { DashboardPage, JoinRoomPage } from '@/app/pages';
import { AppRoutes } from '@/app/router/routePaths';

export const DashboardRoutes: RouteObject[] = [
  {
    path: AppRoutes.DASHBOARD,
    element: <DashboardPage />
  },
  {
    path: AppRoutes.JOIN_ROOM,
    element: <JoinRoomPage />
  }
];
