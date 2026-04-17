import type { SxProps, Theme } from '@mui/material/styles';

export const registerPageCardSx: SxProps<Theme> = {
  maxWidth: 460,
  mx: 'auto'
};

export const registerPageOptionalFieldsSx: SxProps<Theme> = {
  columnGap: 1.5,
  display: 'grid',
  gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, minmax(0, 1fr))' },
  rowGap: 2.5
};
