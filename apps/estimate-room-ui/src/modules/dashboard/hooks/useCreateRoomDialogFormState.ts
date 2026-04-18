import type { UseFormReturn } from 'react-hook-form';

import {
  DashboardCreateRoomLimits
} from '../constants';
import type { DashboardCreateRoomFormValues } from '../types';
import { dashboardDeckPresets } from '../utils';

export const useCreateRoomDialogFormState = (
  form: UseFormReturn<DashboardCreateRoomFormValues>
) => {
  const {
    formState: { errors, isSubmitting, isValid },
    register,
    watch
  } = form;
  const createShareLink = watch('createShareLink');
  const inviteTeamId = watch('inviteTeamId');
  const deckKey = watch('deckKey');
  const roomName = watch('name');
  const roomNameMaxLength = DashboardCreateRoomLimits.ROOM_NAME_MAX_LENGTH;
  const isRoomNameLimitReached = roomName.length >= roomNameMaxLength;
  const roomNameHelperText =
    errors.name?.message ??
    (isRoomNameLimitReached
      ? `Room name limit reached (${roomNameMaxLength} characters).`
      : undefined);

  return {
    isSubmitting,
    isValid,
    deckField: {
      error: Boolean(errors.deckKey),
      helperText: errors.deckKey?.message ?? 'Choose a planning scale.',
      options: dashboardDeckPresets,
      registration: register('deckKey', {
        required: 'Select an estimation deck.'
      }),
      value: deckKey
    },
    inviteEmailsField: {
      error: Boolean(errors.inviteEmails),
      helperText:
        errors.inviteEmails?.message ??
        'Separate participant emails with commas, semicolons, or new lines.',
      registration: register('inviteEmails')
    },
    roomNameField: {
      error: Boolean(errors.name) || isRoomNameLimitReached,
      helperText: roomNameHelperText,
      maxLength: roomNameMaxLength,
      registration: register('name', {
        maxLength: {
          message: `Room names can be up to ${roomNameMaxLength} characters.`,
          value: roomNameMaxLength
        },
        required: 'Room name is required.'
      })
    },
    shareLinkField: {
      checked: createShareLink,
      description: createShareLink
        ? 'A join token will be generated for quick room access.'
        : 'Only invited participants will be able to join.',
      registration: register('createShareLink')
    },
    teamField: {
      error: Boolean(errors.inviteTeamId),
      helperText: errors.inviteTeamId?.message,
      registration: register('inviteTeamId'),
      value: inviteTeamId
    }
  };
};

export type CreateRoomDialogFormState = ReturnType<typeof useCreateRoomDialogFormState>;
