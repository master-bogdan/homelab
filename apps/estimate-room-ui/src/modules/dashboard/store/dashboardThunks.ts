import { createAppAsyncThunk } from '@/shared/store';

import type {
  DashboardActiveRoom,
  DashboardCreateRoomFormValues,
  DashboardTeamSummary
} from '../types';
import { extractInviteToken, getDashboardErrorMessage } from '../utils';
import {
  selectActiveSession,
  selectDashboardView,
  selectRecentRooms
} from './dashboardSelectors';
import { dashboardApi } from './dashboardService';

const buildSectionError = (error: unknown, fallback: string) => ({
  message: getDashboardErrorMessage(error, fallback)
});

export const fetchDashboardPage = createAppAsyncThunk(
  'dashboard/fetchDashboardPage',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      const sessions = await dispatch(
        dashboardApi.endpoints.fetchDashboardSessions.initiate(undefined, {
          forceRefetch: true,
          subscribe: false
        })
      ).unwrap();
      const activeSession = selectActiveSession(sessions);

      const [teamsResult, ledgerResult, activeRoomResult] = await Promise.allSettled([
        dispatch(
          dashboardApi.endpoints.fetchDashboardTeams.initiate(undefined, {
            forceRefetch: true,
            subscribe: false
          })
        ).unwrap(),
        dispatch(
          dashboardApi.endpoints.fetchDashboardLedger.initiate(undefined, {
            forceRefetch: true,
            subscribe: false
          })
        ).unwrap(),
        activeSession
          ? dispatch(
              dashboardApi.endpoints.fetchDashboardRoom.initiate(activeSession.id, {
                forceRefetch: true,
                subscribe: false
              })
            ).unwrap()
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
            ? buildSectionError(
                teamsResult.reason,
                'Teams could not be loaded right now.'
              )
            : null,
        view: selectDashboardView({
          activeRoom,
          sessions,
          teams
        })
      };
    } catch (error) {
      return rejectWithValue(
        getDashboardErrorMessage(error, 'Dashboard data could not be loaded right now.')
      );
    }
  }
);

export const fetchCreateRoomTeams = createAppAsyncThunk(
  'dashboard/fetchCreateRoomTeams',
  async (_, { dispatch, rejectWithValue }) => {
    try {
      return await dispatch(
        dashboardApi.endpoints.fetchDashboardTeams.initiate(undefined, {
          forceRefetch: true,
          subscribe: false
        })
      ).unwrap();
    } catch (error) {
      return rejectWithValue(
        getDashboardErrorMessage(error, 'Teams could not be loaded for room creation.')
      );
    }
  }
);

export const submitCreateRoom = createAppAsyncThunk(
  'dashboard/submitCreateRoom',
  async (values: DashboardCreateRoomFormValues, { dispatch, rejectWithValue }) => {
    try {
      return await dispatch(dashboardApi.endpoints.createRoom.initiate(values)).unwrap();
    } catch (error) {
      return rejectWithValue(
        getDashboardErrorMessage(error, 'The room could not be created right now.')
      );
    }
  }
);

export const submitJoinRoom = createAppAsyncThunk(
  'dashboard/submitJoinRoom',
  async (code: string, { dispatch, rejectWithValue }) => {
    try {
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

      return await dispatch(
        dashboardApi.endpoints.acceptInvitation.initiate(inviteToken)
      ).unwrap();
    } catch (error) {
      return rejectWithValue(
        getDashboardErrorMessage(
          error,
          'Invalid or expired room code. Please check and try again.'
        )
      );
    }
  }
);
