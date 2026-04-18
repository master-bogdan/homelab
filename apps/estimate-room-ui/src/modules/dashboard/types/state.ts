import type {
  DashboardActiveRoom,
  DashboardLedger,
  DashboardSession,
  DashboardTeamSummary
} from './models';
import type { DashboardLoadStatus, DashboardView } from './status';

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
