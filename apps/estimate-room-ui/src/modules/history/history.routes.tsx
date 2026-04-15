import type { RouteObject } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';

import { HistoryPage } from './HistoryPage';
import { HistoryRoomPage } from './HistoryRoomPage';

export const historyRoutes: RouteObject[] = [
  {
    path: appRoutes.history,
    element: <HistoryPage />
  },
  {
    path: appRoutes.historyRoom,
    element: <HistoryRoomPage />
  }
];
