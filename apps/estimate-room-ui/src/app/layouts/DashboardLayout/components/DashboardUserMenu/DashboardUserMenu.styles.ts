import type { SxProps, Theme } from '@mui/material/styles';

export const dashboardUserMenuPaperSx: SxProps<Theme> = {
  border: (theme) => `1px solid ${theme.app.borders.ghost}`,
  boxShadow: (theme) => theme.app.effects.ambientShadow,
  mt: 1,
  minWidth: 240
};

export const dashboardUserMenuSummarySx: SxProps<Theme> = {
  px: 2,
  py: 1.5
};
