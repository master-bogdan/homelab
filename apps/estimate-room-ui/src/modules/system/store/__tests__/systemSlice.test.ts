import {
  clearNotifications,
  closeDialog,
  dismissNotification,
  enqueueNotification,
  openDialog,
  systemReducer
} from '../index';

describe('systemSlice', () => {
  it('opens one dialog at a time and keeps success payloads', () => {
    let state = systemReducer(undefined, openDialog({ key: 'dashboardCreateRoom' }));

    expect(state.dialogs.dashboardCreateRoom.isOpen).toBe(true);
    expect(state.dialogs.dashboardJoinRoom.isOpen).toBe(false);

    state = systemReducer(
      state,
      openDialog({
        key: 'dashboardCreateRoomSuccess',
        payload: {
          result: {
            inviteLink: 'http://localhost:3000/join/ROOM-12',
            roomCode: 'ROOM-12',
            roomId: 'room-1',
            roomName: 'Architecture Review',
            skippedRecipients: []
          }
        }
      })
    );

    expect(state.dialogs.dashboardCreateRoom.isOpen).toBe(false);
    expect(state.dialogs.dashboardCreateRoomSuccess.isOpen).toBe(true);
    expect(state.dialogs.dashboardCreateRoomSuccess.payload?.result.roomId).toBe('room-1');

    state = systemReducer(state, closeDialog('dashboardCreateRoomSuccess'));

    expect(state.dialogs.dashboardCreateRoomSuccess.isOpen).toBe(false);
    expect(state.dialogs.dashboardCreateRoomSuccess.payload).toBeNull();
  });

  it('stores and dismisses queued notifications', () => {
    let state = systemReducer(
      undefined,
      enqueueNotification({
        id: 'notification-1',
        message: 'Room created',
        severity: 'success',
        title: 'Success'
      })
    );

    expect(state.notifications).toHaveLength(1);
    expect(state.notifications[0]).toMatchObject({
      id: 'notification-1',
      message: 'Room created',
      severity: 'success',
      title: 'Success'
    });

    state = systemReducer(state, dismissNotification('notification-1'));
    expect(state.notifications).toHaveLength(0);

    state = systemReducer(
      state,
      enqueueNotification({
        id: 'notification-2',
        message: 'Retry failed',
        severity: 'error'
      })
    );
    state = systemReducer(state, clearNotifications());

    expect(state.notifications).toHaveLength(0);
  });
});
