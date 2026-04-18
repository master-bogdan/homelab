import type { SxProps, Theme } from '@mui/material/styles';

export const notFoundPageRootSx: SxProps<Theme> = {
  display: 'grid',
  minHeight: '100vh',
  placeItems: 'center',
  px: 3,
  py: 6
};

export const notFoundPageCardSx: SxProps<Theme> = {
  maxWidth: 560,
  width: '100%'
};
