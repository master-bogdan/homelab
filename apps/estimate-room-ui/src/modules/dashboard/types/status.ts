import type { DASHBOARD_ROOM_TASK_STATUSES } from '../constants';

export type DashboardSessionStatus = 'ACTIVE' | 'EXPIRED' | 'FINISHED';
export type DashboardRoomTaskStatus =
  (typeof DASHBOARD_ROOM_TASK_STATUSES)[keyof typeof DASHBOARD_ROOM_TASK_STATUSES];
export type DashboardInvitationKind = 'ROOM_EMAIL' | 'ROOM_LINK' | 'TEAM_MEMBER';
export type DashboardLoadStatus = 'error' | 'loading' | 'ready';
export type DashboardView = 'active' | 'empty' | 'noActive';
export type DashboardDeckPresetKey =
  | 'fibonacci'
  | 'powerOfTwo'
  | 'simple'
  | 'tShirt';
