import type { RootState } from '@/app/store/store';

import type {
  DashboardActiveRoom,
  DashboardCreateRoomState,
  DashboardJoinRoomState,
  DashboardPageState,
  DashboardSession,
  DashboardTeamSummary,
  DashboardView
} from '../types';

export const selectActiveSession = (sessions: DashboardSession[]) =>
  sessions.find((session) => session.status === 'ACTIVE') ?? null;

export const selectRecentRooms = (sessions: DashboardSession[], limit = 4) => {
  const uniqueSessions = new Map<string, DashboardSession>();

  sessions.forEach((session) => {
    if (!uniqueSessions.has(session.id)) {
      uniqueSessions.set(session.id, session);
    }
  });

  return Array.from(uniqueSessions.values()).slice(0, limit);
};

export const selectDashboardView = ({
  activeRoom,
  sessions,
  teams
}: {
  readonly activeRoom: DashboardActiveRoom | null;
  readonly sessions: DashboardSession[];
  readonly teams: DashboardTeamSummary[];
}): DashboardView => {
  if (activeRoom) {
    return 'active';
  }

  if (sessions.length === 0 && teams.length === 0) {
    return 'empty';
  }

  return 'noActive';
};

export const selectDashboardState = (state: RootState) => state.dashboard;

export const selectDashboardPageState = (state: RootState): DashboardPageState =>
  selectDashboardState(state).page;

export const selectCreateRoomDialogState = (state: RootState): DashboardCreateRoomState =>
  selectDashboardState(state).createRoom;

export const selectJoinRoomDialogState = (state: RootState): DashboardJoinRoomState =>
  selectDashboardState(state).joinRoom;
