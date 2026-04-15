import type {
  DashboardCreateRoomSuccessDialogPayload,
  SystemDialogKey,
  SystemNotification,
  SystemState
} from '../types';
import { systemStateKey } from './system.store';

type SystemStateRoot = {
  readonly [systemStateKey]: SystemState;
};

export const selectSystemState = (state: SystemStateRoot): SystemState =>
  state[systemStateKey];

export const selectSystemNotifications = (state: SystemStateRoot): SystemNotification[] =>
  selectSystemState(state).notifications;

export const selectIsDialogOpen = (state: SystemStateRoot, key: SystemDialogKey) =>
  selectSystemState(state).dialogs[key].isOpen;

export const selectDashboardCreateRoomSuccessPayload = (
  state: SystemStateRoot
): DashboardCreateRoomSuccessDialogPayload | null =>
  selectSystemState(state).dialogs.dashboardCreateRoomSuccess.payload;
