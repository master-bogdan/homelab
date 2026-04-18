import type { SxProps, Theme } from '@mui/material/styles';

export const dashboardPageStateCardSx: SxProps<Theme> = {
  border: (theme) => `1px solid ${theme.app.borders.ghost}`,
  minHeight: 360,
  p: { xs: 3, md: 4 }
};

export const dashboardPageActiveGridSx: SxProps<Theme> = {
  display: 'grid',
  gap: 4,
  gridTemplateColumns: {
    xs: 'minmax(0, 1fr)',
    xl: 'minmax(0, 6fr) minmax(320px, 4fr)'
  }
};

export const dashboardPageNoActiveSx: SxProps<Theme> = {
  minWidth: 0
};

export const dashboardPageSectionsGridSx: SxProps<Theme> = {
  display: 'grid',
  gap: 4,
  gridTemplateColumns: {
    xs: 'minmax(0, 1fr)',
    xl: 'repeat(2, minmax(0, 1fr))'
  }
};
