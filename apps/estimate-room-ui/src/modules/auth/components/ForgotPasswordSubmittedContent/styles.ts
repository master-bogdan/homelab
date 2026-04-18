import type { SxProps, Theme } from '@mui/material/styles';

export const forgotPasswordSubmittedActionsSx: SxProps<Theme> = {
  borderTop: (theme) => `1px solid ${theme.app.borders.ghost}`,
  pt: 3
};
