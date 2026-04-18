import type { SxProps, Theme } from '@mui/material/styles';

export const forgotPasswordResendLinkSx: SxProps<Theme> = {
  alignItems: 'center',
  background: 'none',
  border: 0,
  cursor: 'pointer',
  display: 'inline-flex',
  gap: 0.75,
  p: 0
};

export const forgotPasswordSubmittedActionsSx: SxProps<Theme> = {
  borderTop: (theme) => `1px solid ${theme.app.borders.ghost}`,
  pt: 3
};
