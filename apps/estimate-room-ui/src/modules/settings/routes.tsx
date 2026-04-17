import type { RouteObject } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';

import { SettingsPage } from './SettingsPage';

export const SettingsRoutes: RouteObject[] = [
  {
    path: AppRoutes.SETTINGS,
    element: <SettingsPage />
  }
];
