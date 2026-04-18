import type { SxProps, Theme } from '@mui/material/styles';

export const authActionDividerLineSx: SxProps<Theme> = {
  bgcolor: (theme) => theme.app.borders.ghost,
  flex: 1,
  height: 1
};
