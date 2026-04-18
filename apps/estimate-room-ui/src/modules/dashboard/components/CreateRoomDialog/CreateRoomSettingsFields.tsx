import {
  AppMenuItem,
  AppStack,
  AppTextField
} from '@/shared/components';

import type { DashboardTeamSummary } from '../../types';
import type { CreateRoomDialogFormState } from '../../hooks';
import { createRoomDialogFieldsRowSx } from './styles';

interface CreateRoomSettingsFieldsProps {
  readonly deckField: CreateRoomDialogFormState['deckField'];
  readonly isLoadingTeams: boolean;
  readonly teamField: CreateRoomDialogFormState['teamField'];
  readonly teamOptions: DashboardTeamSummary[];
}

export const CreateRoomSettingsFields = ({
  deckField,
  isLoadingTeams,
  teamField,
  teamOptions
}: CreateRoomSettingsFieldsProps) => (
  <AppStack direction={{ md: 'row' }} spacing={2} sx={createRoomDialogFieldsRowSx}>
    <AppTextField
      error={teamField.error}
      helperText={isLoadingTeams ? 'Loading teams...' : teamField.helperText ?? 'Optional'}
      label="Team"
      select
      value={teamField.value}
      {...teamField.registration}
    >
      <AppMenuItem value="">No team</AppMenuItem>
      {teamOptions.map((team) => (
        <AppMenuItem key={team.id} value={team.id}>
          {team.name}
        </AppMenuItem>
      ))}
    </AppTextField>
    <AppTextField
      error={deckField.error}
      helperText={deckField.helperText}
      label="Estimation Deck"
      select
      value={deckField.value}
      {...deckField.registration}
    >
      {deckField.options.map((preset) => (
        <AppMenuItem key={preset.key} value={preset.key}>
          {preset.label}
        </AppMenuItem>
      ))}
    </AppTextField>
  </AppStack>
);
