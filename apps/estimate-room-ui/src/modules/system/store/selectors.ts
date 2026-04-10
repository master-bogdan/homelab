import type { RootState } from '@/app/store/store';

import type {
  DashboardCreateRoomSuccessDialogPayload,
  SystemDialogKey,
  SystemNotification,
  SystemState
} from '../types';

export const selectSystemState = (state: RootState): SystemState => state.system;

export const selectSystemNotifications = (state: RootState): SystemNotification[] =>
  selectSystemState(state).notifications;

export const selectIsDialogOpen = (state: RootState, key: SystemDialogKey) =>
  selectSystemState(state).dialogs[key].isOpen;

export const selectDashboardCreateRoomSuccessPayload = (
  state: RootState
): DashboardCreateRoomSuccessDialogPayload | null =>
  selectSystemState(state).dialogs.dashboardCreateRoomSuccess.payload;
