import type { RouteObject } from 'react-router-dom';

import { NewRoomPage, RoomDetailsPage } from '@/app/pages';
import { AppRoutes } from '@/app/router/routePaths';

export const RoomsRoutes: RouteObject[] = [
  {
    path: AppRoutes.ROOMS_NEW,
    element: <NewRoomPage />
  },
  {
    path: AppRoutes.ROOM_DETAILS,
    element: <RoomDetailsPage />
  }
];
