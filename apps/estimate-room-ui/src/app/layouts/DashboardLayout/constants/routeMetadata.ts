import { matchPath } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';

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
    path: appRoutes.dashboard
  },
  {
    meta: {
      description: 'Create a new estimation workspace and invite collaborators.',
      title: 'New Room'
    },
    path: appRoutes.roomsNew
  },
  {
    meta: {
      description: 'Continue the active estimation session and track room progress.',
      title: 'Room'
    },
    path: appRoutes.roomDetails
  },
  {
    meta: {
      description: 'Review completed sessions and archived room outcomes.',
      title: 'History'
    },
    path: appRoutes.history
  },
  {
    meta: {
      description: 'Inspect a completed room session and its estimation output.',
      title: 'Session History'
    },
    path: appRoutes.historyRoom
  },
  {
    meta: {
      description: 'Review members, ownership, and team-linked work.',
      title: 'Team'
    },
    path: appRoutes.teamDetails
  },
  {
    meta: {
      description: 'Manage your profile details and account identity.',
      title: 'Profile'
    },
    path: appRoutes.profile
  },
  {
    meta: {
      description: 'Configure dashboard preferences and account defaults.',
      title: 'Settings'
    },
    path: appRoutes.settings
  }
];

export const resolveDashboardLayoutMeta = (pathname: string) =>
  routeMetadata.find(({ end = true, path }) => matchPath({ end, path }, pathname))?.meta ??
  defaultMeta;
