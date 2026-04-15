import type { SxProps, Theme } from '@mui/material/styles';

export const passwordRecommendationsRootSx: SxProps<Theme> = {
  bgcolor: 'secondary.light',
  borderRadius: (theme) => theme.app.radii.sm,
  px: 2,
  py: 1.75
};

export const passwordRecommendationsGridSx: SxProps<Theme> = {
  columnGap: 3,
  display: 'grid',
  gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, minmax(0, 1fr))' },
  rowGap: 1
};
