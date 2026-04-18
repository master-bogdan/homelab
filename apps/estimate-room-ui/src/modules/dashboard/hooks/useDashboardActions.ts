import { useAppDispatch } from '@/shared/hooks';
import { openDialog } from '@/modules/system/store';

import {
  fetchCreateRoomTeams,
  resetCreateRoomDialogState,
  resetJoinRoomDialogState
} from '../store';

export const useDashboardActions = () => {
  const dispatch = useAppDispatch();

  return {
    openCreateRoom: () => {
      dispatch(resetCreateRoomDialogState());
      dispatch(openDialog({ key: 'dashboardCreateRoom' }));
      dispatch(fetchCreateRoomTeams());
    },
    openJoinRoom: () => {
      dispatch(resetJoinRoomDialogState());
      dispatch(openDialog({ key: 'dashboardJoinRoom' }));
    }
  };
};
