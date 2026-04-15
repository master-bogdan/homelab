import type {
  DashboardActiveRoom,
  DashboardCreateRoomState,
  DashboardJoinRoomState,
  DashboardPageState,
  DashboardSession,
  DashboardTeamSummary,
  DashboardState,
  DashboardView
} from '../types';
import { dashboardStateKey } from './dashboard.store';

type DashboardStateRoot = {
  readonly [dashboardStateKey]: DashboardState;
};

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

export const selectDashboardState = (state: DashboardStateRoot) =>
  state[dashboardStateKey];

export const selectDashboardPageState = (state: DashboardStateRoot): DashboardPageState =>
  selectDashboardState(state).page;

export const selectCreateRoomDialogState = (
  state: DashboardStateRoot
): DashboardCreateRoomState =>
  selectDashboardState(state).createRoom;

export const selectJoinRoomDialogState = (
  state: DashboardStateRoot
): DashboardJoinRoomState =>
  selectDashboardState(state).joinRoom;
