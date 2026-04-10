import type { AppDispatch } from '@/app/store/store';

import {
  selectActiveSession,
  selectDashboardView,
  selectRecentRooms
} from '../store';
import type {
  DashboardActiveRoom,
  DashboardCreateRoomFormValues,
  DashboardCreateRoomResult,
  DashboardJoinRoomResult,
  DashboardPageData,
  DashboardTeamSummary
} from '../types';
import {
  extractInviteToken,
  getDashboardErrorMessage
} from '../utils';
import { dashboardApi } from './dashboardApi';

const fetchSessions = (dispatch: AppDispatch) =>
  dispatch(dashboardApi.endpoints.fetchDashboardSessions.initiate(undefined, {
    forceRefetch: true,
    subscribe: false
  })).unwrap();

const fetchLedger = (dispatch: AppDispatch) =>
  dispatch(dashboardApi.endpoints.fetchDashboardLedger.initiate(undefined, {
    forceRefetch: true,
    subscribe: false
  })).unwrap();

const fetchActiveRoom = (dispatch: AppDispatch, roomId: string) =>
  dispatch(dashboardApi.endpoints.fetchDashboardRoom.initiate(roomId, {
    forceRefetch: true,
    subscribe: false
  })).unwrap();

const buildSectionError = (error: unknown, fallback: string) => ({
  message: getDashboardErrorMessage(error, fallback)
});

export const dashboardService = {
  createRoom: async (
    dispatch: AppDispatch,
    values: DashboardCreateRoomFormValues
  ): Promise<DashboardCreateRoomResult> => {
    return dispatch(dashboardApi.endpoints.createRoom.initiate(values)).unwrap();
  },
  fetchDashboardPageData: async (dispatch: AppDispatch): Promise<DashboardPageData> => {
    const sessions = await fetchSessions(dispatch);
    const activeSession = selectActiveSession(sessions);

    const [teamsResult, ledgerResult, activeRoomResult] = await Promise.allSettled([
      dashboardService.fetchTeams(dispatch),
      fetchLedger(dispatch),
      activeSession
        ? fetchActiveRoom(dispatch, activeSession.id)
        : Promise.resolve<DashboardActiveRoom | null>(null)
    ]);

    const teams = teamsResult.status === 'fulfilled' ? teamsResult.value : [];
    const ledger = ledgerResult.status === 'fulfilled' ? ledgerResult.value : null;
    const activeRoom =
      activeRoomResult.status === 'fulfilled' ? activeRoomResult.value : null;

    return {
      activeRoom,
      activeRoomError:
        activeRoomResult.status === 'rejected'
          ? buildSectionError(
              activeRoomResult.reason,
              'The active room could not be refreshed right now.'
            )
          : null,
      ledger,
      ledgerError:
        ledgerResult.status === 'rejected'
          ? buildSectionError(
              ledgerResult.reason,
              'Architect ledger data is temporarily unavailable.'
            )
          : null,
      recentRooms: selectRecentRooms(sessions),
      sessions,
      teams,
      teamsError:
        teamsResult.status === 'rejected'
          ? buildSectionError(teamsResult.reason, 'Teams could not be loaded right now.')
          : null,
      view: selectDashboardView({
        activeRoom,
        sessions,
        teams
      })
    };
  },
  fetchTeams: async (dispatch: AppDispatch): Promise<DashboardTeamSummary[]> =>
    dispatch(dashboardApi.endpoints.fetchDashboardTeams.initiate(undefined, {
      forceRefetch: true,
      subscribe: false
    })).unwrap(),
  joinRoom: async (dispatch: AppDispatch, code: string): Promise<DashboardJoinRoomResult> => {
    const inviteToken = extractInviteToken(code);

    if (!inviteToken) {
      throw new Error('Enter a room code or paste an invite link.');
    }

    const preview = await dispatch(
      dashboardApi.endpoints.previewInvitation.initiate(inviteToken, {
        forceRefetch: true,
        subscribe: false
      })
    ).unwrap();

    if (preview.kind === 'TEAM_MEMBER') {
      throw new Error('This code belongs to a team invitation, not a room session.');
    }

    if (!preview.roomId) {
      throw new Error('This room code does not point to a valid room.');
    }

    return dispatch(dashboardApi.endpoints.acceptInvitation.initiate(inviteToken)).unwrap();
  }
};
