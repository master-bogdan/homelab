import type { RouteObject } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';

import { ProfilePage } from './ProfilePage';

export const profileRoutes: RouteObject[] = [
  {
    path: appRoutes.profile,
    element: <ProfilePage />
  }
];
