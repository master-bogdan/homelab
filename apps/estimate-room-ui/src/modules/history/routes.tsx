import type { RouteObject } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';

import { HistoryPage } from './HistoryPage';
import { HistoryRoomPage } from './HistoryRoomPage';

export const historyRoutes: RouteObject[] = [
  {
    path: AppRoutes.HISTORY,
    element: <HistoryPage />
  },
  {
    path: AppRoutes.HISTORY_ROOM,
    element: <HistoryRoomPage />
  }
];
