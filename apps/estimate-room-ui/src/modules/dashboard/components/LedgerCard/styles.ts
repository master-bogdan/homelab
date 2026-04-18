import type { SxProps, Theme } from '@mui/material/styles';

export const ledgerCardRootSx = (hasContent: boolean): SxProps<Theme> => ({
  border: (theme) => `1px solid ${theme.app.borders.ghost}`,
  borderRadius: (theme) => theme.app.radii.lg,
  minHeight: 260,
  p: hasContent ? 3.5 : 3
});

export const ledgerCardBadgeSx: SxProps<Theme> = {
  alignItems: 'center',
  bgcolor: (theme) => theme.app.stateLayers.secondaryPanel,
  borderRadius: (theme) => theme.app.radii.md,
  color: 'primary.main',
  display: 'flex',
  height: 64,
  justifyContent: 'center',
  width: 64
};

export const ledgerCardProgressBarSx: SxProps<Theme> = {
  borderRadius: (theme) => theme.app.radii.pill,
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
