import type { SxProps, Theme } from '@mui/material/styles';

export const joinRoomDialogHintSx: SxProps<Theme> = {
  alignItems: 'flex-start',
  bgcolor: (theme) => theme.app.surfaces.section,
  borderRadius: (theme) => theme.app.radii.md,
  display: 'flex',
  gap: 1.5,
  px: 2,
  py: 1.5
};

export const joinRoomDialogHintIconSx: SxProps<Theme> = {
  mt: 0.25
};
