import type { RouteObject } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';

import { TeamDetailsPage } from './TeamDetailsPage';

export const teamsRoutes: RouteObject[] = [
  {
    path: appRoutes.teamDetails,
    element: <TeamDetailsPage />
  }
];
