import type { RouteObject } from 'react-router-dom';

import { ProfilePage } from '@/app/pages';
import { AppRoutes } from '@/app/router/routePaths';

export const ProfileRoutes: RouteObject[] = [
  {
    path: AppRoutes.PROFILE,
    element: <ProfilePage />
  }
];
