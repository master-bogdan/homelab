import RocketLaunchRoundedIcon from '@mui/icons-material/RocketLaunchRounded';
import type { UseFormReturn } from 'react-hook-form';

import {
  AppAlert,
  AppBox,
  AppDialog,
  AppFormControlLabel,
  AppMenuItem,
  AppStack,
  AppSwitch,
  AppTextField,
  AppTypography
} from '@/shared/ui';

import type {
  DashboardCreateRoomFormValues,
  DashboardTeamSummary
} from '../../types';
import { dashboardDeckPresets } from '../../utils';

import {
  createRoomDialogFieldsRowSx,
  createRoomDialogShareLinkPanelSx,
  createRoomDialogSwitchLabelSx
} from './styles';

const roomNameMaxLength = 100;

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
  const {
    formState: { errors, isSubmitting, isValid },
    register,
    watch
  } = form;
  const createShareLink = watch('createShareLink');
  const inviteTeamId = watch('inviteTeamId');
  const deckKey = watch('deckKey');
  const roomName = watch('name');
  const isRoomNameLimitReached = roomName.length >= roomNameMaxLength;
  const roomNameHelperText = errors.name?.message
    ?? (isRoomNameLimitReached
      ? `Room name limit reached (${roomNameMaxLength} characters).`
      : undefined);

  return (
    <AppDialog
      cancelDisabled={isSubmitting}
      confirmDisabled={!isValid}
      confirmLabel="Create Room"
      confirmLoading={isSubmitting}
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
          {submitErrorMessage ? <AppAlert severity="error">{submitErrorMessage}</AppAlert> : null}
          {teamErrorMessage ? (
            <AppAlert severity="warning">
              {teamErrorMessage} You can still create a room without linking a team.
            </AppAlert>
          ) : null}
          <AppTextField
            autoFocus
            error={Boolean(errors.name) || isRoomNameLimitReached}
            helperText={roomNameHelperText}
            inputProps={{ maxLength: roomNameMaxLength }}
            InputLabelProps={{ shrink: true }}
            label="Room Name"
            placeholder="Q4 Infrastructure Scaling"
            {...register('name', {
              maxLength: {
                message: `Room names can be up to ${roomNameMaxLength} characters.`,
                value: roomNameMaxLength
              },
              required: 'Room name is required.'
            })}
          />
          <AppStack direction={{ md: 'row' }} spacing={2} sx={createRoomDialogFieldsRowSx}>
            <AppTextField
              error={Boolean(errors.inviteTeamId)}
              helperText={
                isLoadingTeams ? 'Loading teams...' : errors.inviteTeamId?.message ?? 'Optional'
              }
              label="Team"
              select
              value={inviteTeamId}
              {...register('inviteTeamId')}
            >
              <AppMenuItem value="">No team</AppMenuItem>
              {teamOptions.map((team) => (
                <AppMenuItem key={team.id} value={team.id}>
                  {team.name}
                </AppMenuItem>
              ))}
            </AppTextField>
            <AppTextField
              error={Boolean(errors.deckKey)}
              helperText={errors.deckKey?.message ?? 'Choose a planning scale.'}
              label="Estimation Deck"
              select
              value={deckKey}
              {...register('deckKey', {
                required: 'Select an estimation deck.'
              })}
            >
              {dashboardDeckPresets.map((preset) => (
                <AppMenuItem key={preset.key} value={preset.key}>
                  {preset.label}
                </AppMenuItem>
              ))}
            </AppTextField>
          </AppStack>
          <AppTextField
            error={Boolean(errors.inviteEmails)}
            helperText={
              errors.inviteEmails?.message
              ?? 'Separate participant emails with commas, semicolons, or new lines.'
            }
            label="Invite Participants"
            minRows={3}
            multiline
            placeholder="engineer@company.com, architect@company.com"
            {...register('inviteEmails')}
          />
          <AppBox sx={createRoomDialogShareLinkPanelSx}>
            <AppStack spacing={0.5}>
              <AppTypography variant="subtitle2">Public share link</AppTypography>
              <AppTypography color="text.secondary" variant="caption">
                {createShareLink
                  ? 'A join token will be generated for quick room access.'
                  : 'Only invited participants will be able to join.'}
              </AppTypography>
            </AppStack>
            <AppFormControlLabel
              control={<AppSwitch {...register('createShareLink')} checked={createShareLink} />}
              label=""
              sx={createRoomDialogSwitchLabelSx}
            />
          </AppBox>
        </AppStack>
      </AppBox>
    </AppDialog>
  );
};
