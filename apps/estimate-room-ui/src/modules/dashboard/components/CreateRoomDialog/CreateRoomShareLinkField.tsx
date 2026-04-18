import {
  AppBox,
  AppFormControlLabel,
  AppStack,
  AppSwitch,
  AppTypography
} from '@/shared/components';

import type { CreateRoomDialogFormState } from '../../hooks';
import {
  createRoomDialogShareLinkPanelSx,
  createRoomDialogSwitchLabelSx
} from './styles';

interface CreateRoomShareLinkFieldProps {
  readonly field: CreateRoomDialogFormState['shareLinkField'];
}

export const CreateRoomShareLinkField = ({ field }: CreateRoomShareLinkFieldProps) => (
  <AppBox sx={createRoomDialogShareLinkPanelSx}>
    <AppStack spacing={0.5}>
      <AppTypography variant="subtitle2">Public share link</AppTypography>
      <AppTypography color="text.secondary" variant="caption">
        {field.description}
      </AppTypography>
    </AppStack>
    <AppFormControlLabel
      control={<AppSwitch {...field.registration} checked={field.checked} />}
      label=""
      sx={createRoomDialogSwitchLabelSx}
    />
  </AppBox>
);
