import type { RouteObject } from 'react-router-dom';

import { TeamDetailsPage } from '@/app/pages';
import { AppRoutes } from '@/app/router/routePaths';

export const TeamsRoutes: RouteObject[] = [
  {
    path: AppRoutes.TEAM_DETAILS,
    element: <TeamDetailsPage />
  }
];
