import type { DashboardCreateRoomResult } from '@/modules/dashboard/types';
import type { ThemeMode } from '@/shared/types';

export type SystemDialogKey =
  | 'dashboardCreateRoom'
  | 'dashboardCreateRoomSuccess'
  | 'dashboardJoinRoom';

export type SystemNotificationSeverity = 'error' | 'info' | 'success' | 'warning';

export interface SystemDialogEntry<TPayload> {
  isOpen: boolean;
  payload: TPayload | null;
}

export interface SystemNotification {
  readonly id: string;
  readonly message: string;
  readonly severity: SystemNotificationSeverity;
  readonly title: string | null;
}

export interface DashboardCreateRoomSuccessDialogPayload {
  readonly result: DashboardCreateRoomResult;
}

export interface SystemState {
  dialogs: {
    dashboardCreateRoom: SystemDialogEntry<null>;
    dashboardCreateRoomSuccess: SystemDialogEntry<DashboardCreateRoomSuccessDialogPayload>;
    dashboardJoinRoom: SystemDialogEntry<null>;
  };
  notifications: SystemNotification[];
  ui: SystemUiState;
}

export interface SystemUiState {
  readonly sidebarOpen: boolean;
  readonly themeMode: ThemeMode;
}

export type OpenSystemDialogPayload =
  | { readonly key: 'dashboardCreateRoom' }
  | {
      readonly key: 'dashboardCreateRoomSuccess';
      readonly payload: DashboardCreateRoomSuccessDialogPayload;
    }
  | { readonly key: 'dashboardJoinRoom' };

export interface EnqueueSystemNotificationPayload {
  readonly id?: string;
  readonly message: string;
  readonly severity: SystemNotificationSeverity;
  readonly title?: string | null;
}
