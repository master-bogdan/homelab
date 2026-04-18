import type { RouteObject } from 'react-router-dom';

import { HistoryPage, HistoryRoomPage } from '@/app/pages';
import { AppRoutes } from '@/app/router/routePaths';

export const HistoryRoutes: RouteObject[] = [
  {
    path: AppRoutes.HISTORY,
    element: <HistoryPage />
  },
  {
    path: AppRoutes.HISTORY_ROOM,
    element: <HistoryRoomPage />
  }
];
