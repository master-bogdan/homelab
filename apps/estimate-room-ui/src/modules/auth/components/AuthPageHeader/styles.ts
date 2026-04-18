import type { SxProps, Theme } from '@mui/material/styles';

export const authPageHeaderRootSx: SxProps<Theme> = {
  mb: 4.5,
  textAlign: 'center'
};

export const authPageHeaderIconSx: SxProps<Theme> = {
  alignItems: 'center',
  bgcolor: 'secondary.light',
  border: (theme) => `1px solid ${theme.app.borders.ghost}`,
  borderRadius: (theme) => theme.app.radii.lg,
  display: 'inline-flex',
  height: 52,
  justifyContent: 'center',
  width: 52
};
