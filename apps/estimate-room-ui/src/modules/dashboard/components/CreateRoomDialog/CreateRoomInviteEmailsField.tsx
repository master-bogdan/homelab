import { AppTextField } from '@/shared/components';

import type { CreateRoomDialogFormState } from '../../hooks';

interface CreateRoomInviteEmailsFieldProps {
  readonly field: CreateRoomDialogFormState['inviteEmailsField'];
}

export const CreateRoomInviteEmailsField = ({
  field
}: CreateRoomInviteEmailsFieldProps) => (
  <AppTextField
    error={field.error}
    helperText={field.helperText}
    label="Invite Participants"
    minRows={3}
    multiline
    placeholder="engineer@company.com, architect@company.com"
    {...field.registration}
  />
);
