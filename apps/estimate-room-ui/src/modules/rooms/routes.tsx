import type { RouteObject } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';

import { NewRoomPage } from './NewRoomPage';
import { RoomDetailsPage } from './RoomDetailsPage';

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
