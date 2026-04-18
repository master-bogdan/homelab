import { matchPath } from 'react-router-dom';

import { AppRoutes } from '@/app/router/routePaths';

export interface DashboardLayoutMeta {
  readonly description: string;
  readonly title: string;
}

const defaultMeta: DashboardLayoutMeta = {
  description: 'Manage architectural ledgers and collaborative rooms.',
  title: 'Dashboard'
};

const routeMetadata: Array<{
  readonly end?: boolean;
  readonly meta: DashboardLayoutMeta;
  readonly path: string;
}> = [
  {
    meta: defaultMeta,
    path: AppRoutes.DASHBOARD
  },
  {
    meta: {
      description: 'Create a new estimation workspace and invite collaborators.',
      title: 'New Room'
    },
    path: AppRoutes.ROOMS_NEW
  },
  {
    meta: {
      description: 'Continue the active estimation session and track room progress.',
      title: 'Room'
    },
    path: AppRoutes.ROOM_DETAILS
  },
  {
    meta: {
      description: 'Review completed sessions and archived room outcomes.',
      title: 'History'
    },
    path: AppRoutes.HISTORY
  },
  {
    meta: {
      description: 'Inspect a completed room session and its estimation output.',
      title: 'Session History'
    },
    path: AppRoutes.HISTORY_ROOM
  },
  {
    meta: {
      description: 'Review members, ownership, and team-linked work.',
      title: 'Team'
    },
    path: AppRoutes.TEAM_DETAILS
  },
  {
    meta: {
      description: 'Manage your profile details and account identity.',
      title: 'Profile'
    },
    path: AppRoutes.PROFILE
  },
  {
    meta: {
      description: 'Configure dashboard preferences and account defaults.',
      title: 'Settings'
    },
    path: AppRoutes.SETTINGS
  }
];

export const resolveDashboardLayoutMeta = (pathname: string) =>
  routeMetadata.find(({ end = true, path }) => matchPath({ end, path }, pathname))?.meta ??
  defaultMeta;
