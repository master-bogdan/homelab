export type DashboardSessionStatus = 'ACTIVE' | 'EXPIRED' | 'FINISHED';
export type DashboardInvitationKind = 'ROOM_EMAIL' | 'ROOM_LINK' | 'TEAM_MEMBER';
export type DashboardLoadStatus = 'error' | 'loading' | 'ready';
export type DashboardView = 'active' | 'empty' | 'noActive';
export type DashboardDeckPresetKey =
  | 'fibonacci'
  | 'powerOfTwo'
  | 'simple'
  | 'tShirt';

export interface DashboardSessionDto {
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

export interface DashboardSessionListResponseDto {
  readonly items: DashboardSessionDto[];
  readonly page: number;
  readonly pageSize: number;
  readonly total: number;
}

export interface DashboardTeamSummaryDto {
  readonly createdAt: string;
  readonly name: string;
  readonly ownerUserId: string;
  readonly teamId: string;
}

export interface DashboardGamificationDto {
  readonly achievements: DashboardAchievementDto[];
  readonly stats: DashboardGamificationStatsDto;
}

export interface DashboardAchievementDto {
  readonly key: string;
  readonly level: number;
  readonly unlockedAt: string;
}

export interface DashboardGamificationStatsDto {
  readonly level: number;
  readonly nextLevelXp: number;
  readonly sessionsAdmined: number;
  readonly sessionsParticipated: number;
  readonly tasksEstimated: number;
  readonly xp: number;
}

export interface DashboardRoomDto {
  readonly adminUserId: string;
  readonly code: string;
  readonly createdAt: string;
  readonly deck: DashboardDeckDto;
  readonly finishedAt?: string | null;
  readonly lastActivityAt: string;
  readonly name: string;
  readonly participants?: DashboardRoomParticipantDto[];
  readonly roomId: string;
  readonly status: string;
  readonly tasks?: DashboardRoomTaskDto[];
  readonly teamId?: string | null;
}

export interface DashboardDeckDto {
  readonly kind: string;
  readonly name: string;
  readonly values: string[];
}

export interface DashboardRoomParticipantDto {
  readonly guestName?: string | null;
  readonly joinedAt: string;
  readonly leftAt?: string | null;
  readonly role: string;
  readonly roomId: string;
  readonly roomParticipantId: string;
  readonly user?: DashboardRoomUserDto | null;
  readonly userId?: string | null;
}

export interface DashboardRoomUserDto {
  readonly avatarUrl?: string | null;
  readonly displayName: string;
  readonly email?: string | null;
  readonly userId: string;
}

export interface DashboardRoomTaskDto {
  readonly createdAt: string;
  readonly description?: string | null;
  readonly externalKey?: string | null;
  readonly finalEstimateValue?: string | null;
  readonly isActive: boolean;
  readonly roomId: string;
  readonly status: string;
  readonly taskId: string;
  readonly title: string;
  readonly updatedAt: string;
}

export interface DashboardInvitationPreviewDto {
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

export interface DashboardInvitationWithTokenDto extends DashboardInvitationPreviewDto {
  readonly token: string;
}

export interface DashboardCreateRoomSkippedRecipientDto {
  readonly email?: string | null;
  readonly reason: string;
  readonly userId?: string | null;
}

export interface DashboardCreateRoomResponseDto {
  readonly emailInvites?: DashboardInvitationWithTokenDto[];
  readonly inviteToken?: string;
  readonly room: DashboardRoomDto;
  readonly shareLink?: DashboardInvitationWithTokenDto;
  readonly skippedRecipients?: DashboardCreateRoomSkippedRecipientDto[];
}

export interface DashboardJoinRoomResponseDto {
  readonly participant?: DashboardRoomParticipantDto;
  readonly room?: DashboardRoomDto;
}

export interface DashboardSession {
  readonly approxDurationSeconds: number;
  readonly createdAt: string;
  readonly estimatedTasksCount: number;
  readonly finishedAt: string | null;
  readonly id: string;
  readonly lastActivityAt: string;
  readonly name: string;
  readonly participantsCount: number;
  readonly role: string;
  readonly status: DashboardSessionStatus;
  readonly tasksCount: number;
  readonly teamId: string | null;
}

export interface DashboardTeamSummary {
  readonly createdAt: string;
  readonly id: string;
  readonly name: string;
  readonly ownerUserId: string;
}

export interface DashboardLedger {
  readonly achievements: DashboardAchievement[];
  readonly currentLevelXp: number;
  readonly level: number;
  readonly nextLevelXp: number;
  readonly sessionsAdmined: number;
  readonly sessionsParticipated: number;
  readonly tasksEstimated: number;
  readonly xpProgressPercentage: number;
}

export interface DashboardAchievement {
  readonly key: string;
  readonly level: number;
  readonly unlockedAt: string;
}

export interface DashboardRoomParticipant {
  readonly avatarUrl: string | null;
  readonly displayName: string;
  readonly id: string;
  readonly role: string;
}

export interface DashboardActiveRoom {
  readonly code: string;
  readonly currentTaskStatus: string | null;
  readonly currentTaskTitle: string | null;
  readonly estimatedTasksCount: number;
  readonly id: string;
  readonly lastActivityAt: string;
  readonly name: string;
  readonly participants: DashboardRoomParticipant[];
  readonly status: string;
  readonly tasksCount: number;
  readonly teamId: string | null;
}

export interface DashboardSectionErrorState {
  readonly message: string;
}

export interface DashboardPageData {
  readonly activeRoom: DashboardActiveRoom | null;
  readonly activeRoomError: DashboardSectionErrorState | null;
  readonly ledger: DashboardLedger | null;
  readonly ledgerError: DashboardSectionErrorState | null;
  readonly recentRooms: DashboardSession[];
  readonly sessions: DashboardSession[];
  readonly teams: DashboardTeamSummary[];
  readonly teamsError: DashboardSectionErrorState | null;
  readonly view: DashboardView;
}

export interface DashboardPageState {
  readonly data: DashboardPageData | null;
  readonly errorMessage: string | null;
  readonly status: DashboardLoadStatus;
}

export interface DashboardCreateRoomState {
  readonly isLoadingTeams: boolean;
  readonly submitErrorMessage: string | null;
  readonly teamErrorMessage: string | null;
  readonly teamOptions: DashboardTeamSummary[];
}

export interface DashboardJoinRoomState {
  readonly errorMessage: string | null;
}

export interface DashboardState {
  readonly createRoom: DashboardCreateRoomState;
  readonly joinRoom: DashboardJoinRoomState;
  readonly page: DashboardPageState;
}

export interface DashboardCreateRoomFormValues {
  readonly createShareLink: boolean;
  readonly deckKey: DashboardDeckPresetKey;
  readonly inviteEmails: string;
  readonly inviteTeamId: string;
  readonly name: string;
}

export interface DashboardCreateRoomResult {
  readonly inviteLink: string;
  readonly roomCode: string;
  readonly roomId: string;
  readonly roomName: string;
  readonly skippedRecipients: DashboardCreateRoomSkippedRecipientDto[];
}

export interface DashboardJoinRoomFormValues {
  readonly code: string;
}

export interface DashboardJoinRoomResult {
  readonly roomId: string;
  readonly roomName: string;
}

export interface DashboardDeckPreset {
  readonly description: string;
  readonly key: DashboardDeckPresetKey;
  readonly label: string;
  readonly deck: {
    readonly kind: string;
    readonly name: string;
    readonly values: string[];
  };
}
