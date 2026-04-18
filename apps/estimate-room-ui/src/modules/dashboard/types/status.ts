import type { DashboardLoadStatuses } from '../constants/loadStatuses';
import type { DashboardRoomStatuses } from '../constants/roomStatuses';
import type { DashboardRoomTaskStatuses } from '../constants/taskStatuses';

export type DashboardRoomStatus =
  (typeof DashboardRoomStatuses)[keyof typeof DashboardRoomStatuses];
export type DashboardSessionStatus = DashboardRoomStatus;
export type DashboardRoomTaskStatus =
  (typeof DashboardRoomTaskStatuses)[keyof typeof DashboardRoomTaskStatuses];
export type DashboardInvitationKind = 'ROOM_EMAIL' | 'ROOM_LINK' | 'TEAM_MEMBER';
export type DashboardLoadStatus =
  (typeof DashboardLoadStatuses)[keyof typeof DashboardLoadStatuses];
export type DashboardView = 'active' | 'empty' | 'noActive';
export type DashboardDeckPresetKey =
  | 'fibonacci'
  | 'powerOfTwo'
  | 'simple'
  | 'tShirt';
