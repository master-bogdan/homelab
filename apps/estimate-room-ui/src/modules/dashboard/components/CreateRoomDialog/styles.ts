import type { SxProps, Theme } from '@mui/material/styles';

export const createRoomDialogFieldsRowSx: SxProps<Theme> = {};

export const createRoomDialogShareLinkPanelSx: SxProps<Theme> = {
  alignItems: 'center',
  bgcolor: (theme) => theme.app.stateLayers.secondaryPanel,
  borderRadius: (theme) => theme.app.radii.md,
  display: 'flex',
  justifyContent: 'space-between',
  px: 2,
  py: 1.5
};

export const createRoomDialogSwitchLabelSx: SxProps<Theme> = {
  m: 0
};
