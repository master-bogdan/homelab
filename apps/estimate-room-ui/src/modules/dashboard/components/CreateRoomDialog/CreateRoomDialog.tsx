import RocketLaunchRoundedIcon from '@mui/icons-material/RocketLaunchRounded';
import type { UseFormReturn } from 'react-hook-form';

import {
  AppBox,
  AppDialog,
  AppStack
} from '@/shared/components';

import type {
  DashboardCreateRoomFormValues,
  DashboardTeamSummary
} from '../../types';
import { useCreateRoomDialogFormState } from '../../hooks';

import { CreateRoomDialogAlerts } from './CreateRoomDialogAlerts';
import { CreateRoomInviteEmailsField } from './CreateRoomInviteEmailsField';
import { CreateRoomNameField } from './CreateRoomNameField';
import { CreateRoomSettingsFields } from './CreateRoomSettingsFields';
import { CreateRoomShareLinkField } from './CreateRoomShareLinkField';

export interface CreateRoomDialogProps {
  readonly form: UseFormReturn<DashboardCreateRoomFormValues>;
  readonly isLoadingTeams: boolean;
  readonly open: boolean;
  readonly onClose: () => void;
  readonly onSubmit: () => void;
  readonly submitErrorMessage: string | null;
  readonly teamErrorMessage: string | null;
  readonly teamOptions: DashboardTeamSummary[];
}

export const CreateRoomDialog = ({
  form,
  isLoadingTeams,
  open,
  onClose,
  onSubmit,
  submitErrorMessage,
  teamErrorMessage,
  teamOptions
}: CreateRoomDialogProps) => {
  const formState = useCreateRoomDialogFormState(form);

  return (
    <AppDialog
      cancelDisabled={formState.isSubmitting}
      confirmDisabled={!formState.isValid}
      confirmLabel="Create Room"
      confirmLoading={formState.isSubmitting}
      confirmLoadingText="Creating..."
      confirmStartIcon={<RocketLaunchRoundedIcon />}
      maxWidth="sm"
      onCancel={onClose}
      onClose={onClose}
      onConfirm={onSubmit}
      open={open}
      title="Create New Room"
    >
      <AppBox component="form" noValidate onSubmit={onSubmit}>
        <AppStack spacing={2.5}>
          <CreateRoomDialogAlerts
            submitErrorMessage={submitErrorMessage}
            teamErrorMessage={teamErrorMessage}
          />
          <CreateRoomNameField field={formState.roomNameField} />
          <CreateRoomSettingsFields
            deckField={formState.deckField}
            isLoadingTeams={isLoadingTeams}
            teamField={formState.teamField}
            teamOptions={teamOptions}
          />
          <CreateRoomInviteEmailsField field={formState.inviteEmailsField} />
          <CreateRoomShareLinkField field={formState.shareLinkField} />
        </AppStack>
      </AppBox>
    </AppDialog>
  );
};
