import { createSlice, type PayloadAction } from '@reduxjs/toolkit';

import type {
  EnqueueSystemNotificationPayload,
  OpenSystemDialogPayload,
  SystemDialogEntry,
  SystemDialogKey,
  SystemNotification,
  SystemState
} from '../types';

const createDialogEntry = <TPayload>(): SystemDialogEntry<TPayload> => ({
  isOpen: false,
  payload: null
});

const initialState: SystemState = {
  dialogs: {
    dashboardCreateRoom: createDialogEntry(),
    dashboardCreateRoomSuccess: createDialogEntry(),
    dashboardJoinRoom: createDialogEntry()
  },
  notifications: []
};

const closeAllDialogs = (state: SystemState) => {
  state.dialogs.dashboardCreateRoom.isOpen = false;
  state.dialogs.dashboardCreateRoom.payload = null;
  state.dialogs.dashboardCreateRoomSuccess.isOpen = false;
  state.dialogs.dashboardCreateRoomSuccess.payload = null;
  state.dialogs.dashboardJoinRoom.isOpen = false;
  state.dialogs.dashboardJoinRoom.payload = null;
};

const systemSlice = createSlice({
  name: 'system',
  initialState,
  reducers: {
    clearNotifications: (state) => {
      state.notifications = [];
    },
    closeDialog: (state, action: PayloadAction<SystemDialogKey>) => {
      const dialog = state.dialogs[action.payload];
      dialog.isOpen = false;
      dialog.payload = null;
    },
    dismissNotification: (state, action: PayloadAction<string>) => {
      state.notifications = state.notifications.filter(
        (notification) => notification.id !== action.payload
      );
    },
    enqueueNotification: {
      reducer: (state, action: PayloadAction<SystemNotification>) => {
        state.notifications.push(action.payload);
      },
      prepare: (payload: EnqueueSystemNotificationPayload) => ({
        payload: {
          id: payload.id ?? crypto.randomUUID(),
          message: payload.message,
          severity: payload.severity,
          title: payload.title ?? null
        } satisfies SystemNotification
      })
    },
    openDialog: (state, action: PayloadAction<OpenSystemDialogPayload>) => {
      closeAllDialogs(state);

      const dialog = state.dialogs[action.payload.key];
      dialog.isOpen = true;

      if ('payload' in action.payload) {
        dialog.payload = action.payload.payload;
      }
    }
  }
});

export const {
  clearNotifications,
  closeDialog,
  dismissNotification,
  enqueueNotification,
  openDialog
} = systemSlice.actions;
export const systemReducer = systemSlice.reducer;
