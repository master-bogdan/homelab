import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';

import { useAppDispatch, useAppSelector } from '@/shared/store';
import {
  closeDialog,
  openDialog,
  selectDashboardCreateRoomSuccessPayload,
  selectIsDialogOpen
} from '@/modules/system/store';
import { appRoutes } from '@/shared/constants/routes';

import {
  fetchDashboardPage,
  resetCreateRoomDialogState,
  selectCreateRoomDialogState,
  submitCreateRoom
} from '../store';
import type { DashboardCreateRoomFormValues } from '../types';

const defaultValues: DashboardCreateRoomFormValues = {
  createShareLink: true,
  deckKey: 'fibonacci',
  inviteEmails: '',
  inviteTeamId: '',
  name: ''
};

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
    defaultValues,
    mode: 'onChange'
  });
  const result = successPayload?.result ?? null;

  useEffect(() => {
    if (!isOpen) {
      return;
    }

    form.reset(defaultValues);
    form.clearErrors();
  }, [form, isOpen]);

  const close = () => {
    dispatch(resetCreateRoomDialogState());
    dispatch(closeDialog('dashboardCreateRoom'));
    form.reset(defaultValues);
    form.clearErrors();
  };

  const closeResult = () => {
    dispatch(closeDialog('dashboardCreateRoomSuccess'));
  };

  const openCreatedRoom = () => {
    if (!result) {
      return;
    }

    navigate(appRoutes.roomDetailsPath(result.roomId));
    dispatch(closeDialog('dashboardCreateRoomSuccess'));
  };

  const onSubmit = form.handleSubmit(async (values) => {
    try {
      const nextResult = await dispatch(submitCreateRoom(values)).unwrap();
      dispatch(
        openDialog({
          key: 'dashboardCreateRoomSuccess',
          payload: {
            result: nextResult
          }
        })
      );
      void dispatch(fetchDashboardPage());
    } catch {}
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
