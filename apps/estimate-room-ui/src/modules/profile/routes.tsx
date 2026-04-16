import type { RouteObject } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';

import { ProfilePage } from './ProfilePage';

export const profileRoutes: RouteObject[] = [
  {
    path: AppRoutes.PROFILE,
    element: <ProfilePage />
  }
];
