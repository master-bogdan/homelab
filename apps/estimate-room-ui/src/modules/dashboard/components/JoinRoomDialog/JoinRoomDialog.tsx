import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import type { UseFormReturn } from 'react-hook-form';

import { AppAlert, AppBox, AppDialog, AppStack, AppTextField, AppTypography } from '@/shared/ui';

import type { DashboardJoinRoomFormValues } from '../../types';

import { joinRoomDialogHintIconSx, joinRoomDialogHintSx } from './styles';

export interface JoinRoomDialogProps {
  readonly errorMessage: string | null;
  readonly form: UseFormReturn<DashboardJoinRoomFormValues>;
  readonly onClose: () => void;
  readonly onSubmit: () => void;
  readonly open: boolean;
}

export const JoinRoomDialog = ({
  errorMessage,
  form,
  onClose,
  onSubmit,
  open
}: JoinRoomDialogProps) => {
  const {
    formState: { errors, isSubmitting },
    register
  } = form;

  return (
    <AppDialog
      cancelDisabled={isSubmitting}
      confirmLabel="Join Room"
      confirmLoading={isSubmitting}
      confirmLoadingText="Joining..."
      maxWidth="xs"
      onCancel={onClose}
      onClose={onClose}
      onConfirm={onSubmit}
      open={open}
      title="Join Room"
    >
      <AppBox component="form" noValidate onSubmit={onSubmit}>
        <AppStack spacing={2.5}>
          {errorMessage ? <AppAlert severity="error">{errorMessage}</AppAlert> : null}
          <AppTextField
            autoFocus
            error={Boolean(errors.code)}
            helperText={errors.code?.message ?? 'Paste the room invite token or a full invite link.'}
            label="Room Code"
            placeholder="Paste the invite token"
            {...register('code', {
              required: 'Enter a room code or invite link.'
            })}
          />
          <AppBox sx={joinRoomDialogHintSx}>
            <InfoOutlinedIcon color="primary" fontSize="small" sx={joinRoomDialogHintIconSx} />
            <AppTypography color="text.secondary" variant="caption">
              Team invitations are intentionally skipped in this dashboard flow. Only room
              session codes are accepted here.
            </AppTypography>
          </AppBox>
        </AppStack>
      </AppBox>
    </AppDialog>
  );
};
