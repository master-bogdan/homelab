import type { SxProps, Theme } from '@mui/material';

import { APP_DRAWER_WIDTH } from '@/shared/constants/layout';

export const dashboardSidebarDrawerSx: SxProps<Theme> = {
  '& .MuiDrawer-paper': {
    borderRight: 'none',
    boxSizing: 'border-box',
    width: APP_DRAWER_WIDTH
  },
  width: APP_DRAWER_WIDTH
};

export const dashboardSidebarContentSx: SxProps<Theme> = {
  bgcolor: (theme) => theme.app.surfaces.section,
  display: 'flex',
  flexDirection: 'column',
  height: '100%',
  px: 2,
  py: 3
};

export const dashboardSidebarBrandSx: SxProps<Theme> = {
  mb: 5,
  px: 1
};

export const dashboardSidebarBrandIconSx: SxProps<Theme> = {
  alignItems: 'center',
  backgroundImage: (theme) => theme.app.gradients.primary,
  borderRadius: (theme) => theme.app.radii.md,
  color: 'primary.contrastText',
  display: 'flex',
  height: 40,
  justifyContent: 'center',
  width: 40
};

export const dashboardSidebarListSx: SxProps<Theme> = {
  display: 'grid',
  gap: 0.5
};

export const dashboardSidebarSecondaryListSx: SxProps<Theme> = {
  ...dashboardSidebarListSx,
  mt: 'auto'
};

export const dashboardSidebarItemSx: SxProps<Theme> = {
  borderRadius: (theme) => theme.app.radii.md
};

export const dashboardSidebarItemIconSx: SxProps<Theme> = {
  color: 'inherit',
  minWidth: 40
};
