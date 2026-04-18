import type {
  DashboardInvitationKind,
  DashboardRoomStatus,
  DashboardRoomTaskStatus,
  DashboardSessionStatus
} from './status';

export interface DashboardSessionApiResponse {
  readonly approxDurationSeconds: number;
  readonly createdAt: string;
  readonly estimatedTasksCount: number;
  readonly finishedAt?: string | null;
  readonly lastActivityAt: string;
  readonly name: string;
  readonly participantsCount: number;
  readonly role: string;
  readonly roomId: string;
  readonly status: DashboardSessionStatus;
  readonly tasksCount: number;
  readonly teamId?: string | null;
}

export interface DashboardSessionListApiResponse {
  readonly items: DashboardSessionApiResponse[];
  readonly page: number;
  readonly pageSize: number;
  readonly total: number;
}

export interface DashboardTeamSummaryApiResponse {
  readonly createdAt: string;
  readonly name: string;
  readonly ownerUserId: string;
  readonly teamId: string;
}

export interface DashboardGamificationApiResponse {
  readonly achievements: DashboardAchievementApiResponse[];
  readonly stats: DashboardGamificationStatsApiResponse;
}

export interface DashboardAchievementApiResponse {
  readonly key: string;
  readonly level: number;
  readonly unlockedAt: string;
}

export interface DashboardGamificationStatsApiResponse {
  readonly level: number;
  readonly nextLevelXp: number;
  readonly sessionsAdmined: number;
  readonly sessionsParticipated: number;
  readonly tasksEstimated: number;
  readonly xp: number;
}

export interface DashboardRoomApiResponse {
  readonly adminUserId: string;
  readonly code: string;
  readonly createdAt: string;
  readonly deck: DashboardDeckApiResponse;
  readonly finishedAt?: string | null;
  readonly lastActivityAt: string;
  readonly name: string;
  readonly participants?: DashboardRoomParticipantApiResponse[];
  readonly roomId: string;
  readonly status: DashboardRoomStatus;
  readonly tasks?: DashboardRoomTaskApiResponse[];
  readonly teamId?: string | null;
}

export interface DashboardDeckApiResponse {
  readonly kind: string;
  readonly name: string;
  readonly values: string[];
}

export interface DashboardRoomParticipantApiResponse {
  readonly guestName?: string | null;
  readonly joinedAt: string;
  readonly leftAt?: string | null;
  readonly role: string;
  readonly roomId: string;
  readonly roomParticipantId: string;
  readonly user?: DashboardRoomUserApiResponse | null;
  readonly userId?: string | null;
}

export interface DashboardRoomUserApiResponse {
  readonly avatarUrl?: string | null;
  readonly displayName: string;
  readonly email?: string | null;
  readonly userId: string;
}

export interface DashboardRoomTaskApiResponse {
  readonly createdAt: string;
  readonly description?: string | null;
  readonly externalKey?: string | null;
  readonly finalEstimateValue?: string | null;
  readonly isActive: boolean;
  readonly roomId: string;
  readonly status: DashboardRoomTaskStatus;
  readonly taskId: string;
  readonly title: string;
  readonly updatedAt: string;
}

export interface DashboardInvitationPreviewApiResponse {
  readonly acceptedAt?: string | null;
  readonly createdAt: string;
  readonly createdByUserId: string;
  readonly declinedAt?: string | null;
  readonly invitedEmail?: string | null;
  readonly invitedUserId?: string | null;
  readonly invitationId: string;
  readonly kind: DashboardInvitationKind;
  readonly revokedAt?: string | null;
  readonly roomId?: string | null;
  readonly status: string;
  readonly teamId?: string | null;
  readonly updatedAt: string;
}

export interface DashboardInvitationWithTokenApiResponse
  extends DashboardInvitationPreviewApiResponse {
  readonly token: string;
}

export interface DashboardCreateRoomSkippedRecipientApiResponse {
  readonly email?: string | null;
  readonly reason: string;
  readonly userId?: string | null;
}

export interface DashboardCreateRoomApiResponse {
  readonly emailInvites?: DashboardInvitationWithTokenApiResponse[];
  readonly inviteToken?: string;
  readonly room: DashboardRoomApiResponse;
  readonly shareLink?: DashboardInvitationWithTokenApiResponse;
  readonly skippedRecipients?: DashboardCreateRoomSkippedRecipientApiResponse[];
}

export interface DashboardJoinRoomApiResponse {
  readonly participant?: DashboardRoomParticipantApiResponse;
  readonly room?: DashboardRoomApiResponse;
}
