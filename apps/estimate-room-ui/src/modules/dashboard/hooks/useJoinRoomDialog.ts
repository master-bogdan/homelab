import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';

import { useAppDispatch, useAppSelector } from '@/shared/hooks';
import { closeDialog, selectIsDialogOpen } from '@/modules/system/store';
import { AppRoutes } from '@/app/router/routePaths';

import {
  resetJoinRoomDialogState,
  selectJoinRoomDialogState,
  submitJoinRoom
} from '../store';
import type { DashboardJoinRoomFormValues } from '../types';
import { getDashboardErrorMessage } from '../utils';

const defaultValues: DashboardJoinRoomFormValues = {
  code: ''
};

export const useJoinRoomDialog = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const isOpen = useAppSelector((state) => selectIsDialogOpen(state, 'dashboardJoinRoom'));
  const { errorMessage } = useAppSelector(selectJoinRoomDialogState);
  const form = useForm<DashboardJoinRoomFormValues>({
    defaultValues,
    mode: 'onSubmit'
  });

  useEffect(() => {
    if (!isOpen) {
      return;
    }

    form.reset(defaultValues);
    form.clearErrors();
  }, [form, isOpen]);

  const close = () => {
    dispatch(resetJoinRoomDialogState());
    dispatch(closeDialog('dashboardJoinRoom'));
    form.reset(defaultValues);
    form.clearErrors();
  };

  const onSubmit = form.handleSubmit(async (values) => {
    form.clearErrors();

    const result = await dispatch(submitJoinRoom(values.code));

    if (submitJoinRoom.rejected.match(result)) {
      const nextErrorMessage = getDashboardErrorMessage(
        result.payload,
        'Invalid or expired room code. Please check and try again.'
      );

      form.setError('code', {
        message: nextErrorMessage,
        type: 'server'
      });
      return;
    }

    dispatch(resetJoinRoomDialogState());
    dispatch(closeDialog('dashboardJoinRoom'));
    navigate(AppRoutes.ROOM_DETAILS_PATH(result.payload.roomId));
  });

  return {
    close,
    errorMessage,
    form,
    isOpen,
    onSubmit
  };
};
