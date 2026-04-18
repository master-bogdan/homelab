import { AppTextField } from '@/shared/components';

import type { CreateRoomDialogFormState } from '../../hooks';

interface CreateRoomNameFieldProps {
  readonly field: CreateRoomDialogFormState['roomNameField'];
}

export const CreateRoomNameField = ({ field }: CreateRoomNameFieldProps) => (
  <AppTextField
    autoFocus
    error={field.error}
    helperText={field.helperText}
    inputProps={{ maxLength: field.maxLength }}
    InputLabelProps={{ shrink: true }}
    label="Room Name"
    placeholder="Q4 Infrastructure Scaling"
    {...field.registration}
  />
);
