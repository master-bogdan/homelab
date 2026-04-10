import type { SxProps, Theme } from '@mui/material/styles';

export const ledgerCardRootSx = (hasContent: boolean): SxProps<Theme> => ({
  border: (theme) => `1px solid ${theme.app.borders.ghost}`,
  borderRadius: 3,
  minHeight: 260,
  p: hasContent ? 3.5 : 3
});

export const ledgerCardBadgeSx: SxProps<Theme> = {
  alignItems: 'center',
  bgcolor: 'secondary.light',
  borderRadius: 2,
  color: 'primary.main',
  display: 'flex',
  height: 64,
  justifyContent: 'center',
  width: 64
};

export const ledgerCardProgressBarSx: SxProps<Theme> = {
  borderRadius: 999,
  height: 8
};

export const ledgerCardMetricsGridSx: SxProps<Theme> = {
  display: 'grid',
  gap: 2,
  gridTemplateColumns: {
    xs: 'repeat(2, minmax(0, 1fr))',
    sm: 'repeat(4, minmax(0, 1fr))'
  }
};
