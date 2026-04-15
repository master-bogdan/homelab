import type { DashboardRoomTaskStatus, DashboardSessionStatus } from './status';

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
  readonly currentTaskStatus: DashboardRoomTaskStatus | null;
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
