import type { RouteObject } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';

import { SettingsPage } from './SettingsPage';

export const settingsRoutes: RouteObject[] = [
  {
    path: appRoutes.settings,
    element: <SettingsPage />
  }
];
