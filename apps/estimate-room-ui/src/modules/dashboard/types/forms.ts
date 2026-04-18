import type { DashboardDeckPresetKey } from './status';

export interface DashboardCreateRoomFormValues {
  readonly createShareLink: boolean;
  readonly deckKey: DashboardDeckPresetKey;
  readonly inviteEmails: string;
  readonly inviteTeamId: string;
  readonly name: string;
}

export interface DashboardSkippedRecipient {
  readonly email: string | null;
  readonly reason: string;
  readonly userId: string | null;
}

export interface DashboardCreateRoomResult {
  readonly inviteLink: string;
  readonly roomCode: string;
  readonly roomId: string;
  readonly roomName: string;
  readonly skippedRecipients: DashboardSkippedRecipient[];
}

export interface DashboardJoinRoomFormValues {
  readonly code: string;
}

export interface DashboardJoinRoomResult {
  readonly roomId: string;
  readonly roomName: string;
}
