import type { RouteObject } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';

import { TeamDetailsPage } from './TeamDetailsPage';

export const teamsRoutes: RouteObject[] = [
  {
    path: AppRoutes.TEAM_DETAILS,
    element: <TeamDetailsPage />
  }
];
