import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';

import { useAppDispatch, useAppSelector } from '@/shared/hooks';
import {
  closeDialog,
  openDialog,
  selectDashboardCreateRoomSuccessPayload,
  selectIsDialogOpen
} from '@/modules/system/store';
import { AppRoutes } from '@/app/router/routePaths';

import { DashboardCreateRoomDefaultValues } from '../constants';
import {
  resetCreateRoomDialogState,
  selectCreateRoomDialogState,
  submitCreateRoom
} from '../store';
import type { DashboardCreateRoomFormValues } from '../types';

export const useCreateRoomDialog = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const isOpen = useAppSelector((state) => selectIsDialogOpen(state, 'dashboardCreateRoom'));
  const successPayload = useAppSelector(selectDashboardCreateRoomSuccessPayload);
  const {
    isLoadingTeams,
    submitErrorMessage,
    teamErrorMessage,
    teamOptions
  } = useAppSelector(selectCreateRoomDialogState);
  const form = useForm<DashboardCreateRoomFormValues>({
    defaultValues: DashboardCreateRoomDefaultValues,
    mode: 'onChange'
  });
  const result = successPayload?.result ?? null;

  useEffect(() => {
    if (!isOpen) {
      return;
    }

    form.reset(DashboardCreateRoomDefaultValues);
    form.clearErrors();
  }, [form, isOpen]);

  const close = () => {
    dispatch(resetCreateRoomDialogState());
    dispatch(closeDialog('dashboardCreateRoom'));
    form.reset(DashboardCreateRoomDefaultValues);
    form.clearErrors();
  };

  const closeResult = () => {
    dispatch(closeDialog('dashboardCreateRoomSuccess'));
  };

  const openCreatedRoom = () => {
    if (!result) {
      return;
    }

    navigate(AppRoutes.ROOM_DETAILS_PATH(result.roomId));
    dispatch(closeDialog('dashboardCreateRoomSuccess'));
  };

  const onSubmit = form.handleSubmit(async (values) => {
    const submitResult = await dispatch(submitCreateRoom(values));

    if (submitCreateRoom.rejected.match(submitResult)) {
      return;
    }

    dispatch(
      openDialog({
        key: 'dashboardCreateRoomSuccess',
        payload: {
          result: submitResult.payload
        }
      })
    );
  });

  return {
    close,
    closeResult,
    form,
    isLoadingTeams,
    isOpen,
    onSubmit,
    openCreatedRoom,
    result,
    submitErrorMessage,
    teamErrorMessage,
    teamOptions
  };
};
