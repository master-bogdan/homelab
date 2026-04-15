import type { RouteObject } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';

import { NewRoomPage } from './NewRoomPage';
import { RoomDetailsPage } from './RoomDetailsPage';

export const roomsRoutes: RouteObject[] = [
  {
    path: appRoutes.roomsNew,
    element: <NewRoomPage />
  },
  {
    path: appRoutes.roomDetails,
    element: <RoomDetailsPage />
  }
];
