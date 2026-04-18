import { api } from '@/shared/api';

import { DASHBOARD_ROOM_TASK_STATUSES } from '../constants';
import type {
  DashboardActiveRoom,
  DashboardCreateRoomApiResponse,
  DashboardCreateRoomFormValues,
  DashboardCreateRoomResult,
  DashboardCreateRoomSkippedRecipientApiResponse,
  DashboardGamificationApiResponse,
  DashboardInvitationPreviewApiResponse,
  DashboardJoinRoomApiResponse,
  DashboardJoinRoomResult,
  DashboardLedger,
  DashboardRoomApiResponse,
  DashboardRoomParticipant,
  DashboardRoomParticipantApiResponse,
  DashboardSession,
  DashboardSessionApiResponse,
  DashboardSessionListApiResponse,
  DashboardTeamSummary,
  DashboardTeamSummaryApiResponse
} from '../types';
import { buildDashboardInviteLink, getDashboardDeckPreset, parseInviteEmails } from '../utils';

const DASHBOARD_PAGE_SIZE = 20;

const mapSession = (session: DashboardSessionApiResponse): DashboardSession => ({
  approxDurationSeconds: session.approxDurationSeconds,
  createdAt: session.createdAt,
  estimatedTasksCount: session.estimatedTasksCount,
  finishedAt: session.finishedAt ?? null,
  id: session.roomId,
  lastActivityAt: session.lastActivityAt,
  name: session.name,
  participantsCount: session.participantsCount,
  role: session.role,
  status: session.status,
  tasksCount: session.tasksCount,
  teamId: session.teamId ?? null
});

const mapTeam = (team: DashboardTeamSummaryApiResponse): DashboardTeamSummary => ({
  createdAt: team.createdAt,
  id: team.teamId,
  name: team.name,
  ownerUserId: team.ownerUserId
});

const mapLedger = (response: DashboardGamificationApiResponse): DashboardLedger => ({
  achievements: response.achievements.map((achievement) => ({
    key: achievement.key,
    level: achievement.level,
    unlockedAt: achievement.unlockedAt
  })),
  currentLevelXp: response.stats.xp,
  level: response.stats.level,
  nextLevelXp: response.stats.nextLevelXp,
  sessionsAdmined: response.stats.sessionsAdmined,
  sessionsParticipated: response.stats.sessionsParticipated,
  tasksEstimated: response.stats.tasksEstimated,
  xpProgressPercentage:
    response.stats.nextLevelXp <= 0
      ? 100
      : Math.min(100, Math.round((response.stats.xp / response.stats.nextLevelXp) * 100))
});

const mapParticipant = (
  participant: DashboardRoomParticipantApiResponse
): DashboardRoomParticipant => ({
  avatarUrl: participant.user?.avatarUrl ?? null,
  displayName: participant.user?.displayName ?? participant.guestName ?? 'Guest participant',
  id: participant.roomParticipantId,
  role: participant.role
});

const mapActiveRoom = (room: DashboardRoomApiResponse): DashboardActiveRoom => {
  const tasks = room.tasks ?? [];
  const currentTask = tasks.find((task) => task.isActive) ?? null;
  const activeParticipants = (room.participants ?? []).filter((participant) => !participant.leftAt);

  return {
    code: room.code,
    currentTaskStatus: currentTask?.status ?? null,
    currentTaskTitle: currentTask?.title ?? null,
    estimatedTasksCount: tasks.filter(
      (task) =>
        task.status === DASHBOARD_ROOM_TASK_STATUSES.ESTIMATED ||
        task.finalEstimateValue != null
    ).length,
    id: room.roomId,
    lastActivityAt: room.lastActivityAt,
    name: room.name,
    participants: activeParticipants.map(mapParticipant),
    status: room.status,
    tasksCount: tasks.length,
    teamId: room.teamId ?? null
  };
};

const mapCreateRoomResult = (
  response: DashboardCreateRoomApiResponse
): DashboardCreateRoomResult => {
  const shareToken = response.shareLink?.token ?? response.inviteToken;

  if (!shareToken) {
    throw new Error('The room was created, but no invitation token was returned.');
  }

  return {
    inviteLink: buildDashboardInviteLink(shareToken),
    roomCode: shareToken,
    roomId: response.room.roomId,
    roomName: response.room.name,
    skippedRecipients: (response.skippedRecipients ?? []).map(mapSkippedRecipient)
  };
};

const mapSkippedRecipient = (
  recipient: DashboardCreateRoomSkippedRecipientApiResponse
) => ({
  email: recipient.email ?? null,
  reason: recipient.reason,
  userId: recipient.userId ?? null
});

export const dashboardApi = api.injectEndpoints({
  endpoints: (builder) => ({
    acceptInvitation: builder.mutation<DashboardJoinRoomResult, string>({
      query: (token) => ({
        body: {},
        method: 'POST',
        url: `invites/${encodeURIComponent(token)}/accept`
      }),
      transformResponse: (response: DashboardJoinRoomApiResponse) => {
        if (!response.room) {
          throw new Error('The room accepted successfully, but no room details were returned.');
        }

        return {
          roomId: response.room.roomId,
          roomName: response.room.name
        };
      }
    }),
    createRoom: builder.mutation<DashboardCreateRoomResult, DashboardCreateRoomFormValues>({
      query: (values) => {
        const deckPreset = getDashboardDeckPreset(values.deckKey);

        return {
          body: {
            createShareLink: values.createShareLink,
            deck: deckPreset.deck,
            inviteEmails: parseInviteEmails(values.inviteEmails),
            inviteTeamId: values.inviteTeamId || undefined,
            name: values.name.trim()
          },
          method: 'POST',
          url: 'rooms'
        };
      },
      transformResponse: mapCreateRoomResult
    }),
    fetchDashboardLedger: builder.query<DashboardLedger, void>({
      query: () => ({
        url: 'gamification/me'
      }),
      transformResponse: mapLedger
    }),
    fetchDashboardRoom: builder.query<DashboardActiveRoom, string>({
      query: (roomId) => ({
        url: `rooms/${roomId}`
      }),
      transformResponse: mapActiveRoom
    }),
    fetchDashboardSessions: builder.query<DashboardSession[], void>({
      query: () => ({
        params: {
          page: 1,
          pageSize: DASHBOARD_PAGE_SIZE
        },
        url: 'history/me/sessions'
      }),
      transformResponse: (response: DashboardSessionListApiResponse) =>
        response.items.map(mapSession)
    }),
    fetchDashboardTeams: builder.query<DashboardTeamSummary[], void>({
      query: () => ({
        url: 'teams'
      }),
      transformResponse: (response: DashboardTeamSummaryApiResponse[]) => response.map(mapTeam)
    }),
    previewInvitation: builder.query<DashboardInvitationPreviewApiResponse, string>({
      query: (token) => ({
        url: `invites/${encodeURIComponent(token)}`
      })
    })
  }),
  overrideExisting: false
});

export const {
  useCreateRoomMutation,
  useFetchDashboardLedgerQuery,
  useFetchDashboardRoomQuery,
  useFetchDashboardSessionsQuery,
  useFetchDashboardTeamsQuery,
  usePreviewInvitationQuery
} = dashboardApi;
