import type { SxProps, Theme } from '@mui/material';

export const dashboardLayoutRootSx: SxProps<Theme> = {
  bgcolor: 'background.default',
  display: 'flex',
  minHeight: '100vh'
};

export const dashboardLayoutBodySx: SxProps<Theme> = {
  display: 'flex',
  flex: 1,
  flexDirection: 'column',
  minWidth: 0
};

export const dashboardLayoutMainSx: SxProps<Theme> = {
  flex: 1,
  px: { xs: 2, md: 4 },
  py: { xs: 3, md: 4 }
};
