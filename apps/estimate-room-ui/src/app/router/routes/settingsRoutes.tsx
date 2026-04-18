import type { RouteObject } from 'react-router-dom';

import { SettingsPage } from '@/app/pages';
import { AppRoutes } from '@/app/router/routePaths';

export const SettingsRoutes: RouteObject[] = [
  {
    path: AppRoutes.SETTINGS,
    element: <SettingsPage />
  }
];
