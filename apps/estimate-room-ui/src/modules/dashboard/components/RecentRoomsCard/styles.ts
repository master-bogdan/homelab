import type { SxProps, Theme } from '@mui/material/styles';

export const recentRoomsCardActionLinkSx: SxProps<Theme> = {
  textDecoration: 'none'
};

export const recentRoomsCardRootSx = (isEmpty: boolean): SxProps<Theme> => ({
  bgcolor: (theme) => theme.app.surfaces.section,
  borderRadius: (theme) => theme.app.radii.lg,
  minHeight: 320,
  overflow: 'hidden',
  p: isEmpty ? 3 : 1
});

export const recentRoomsCardItemLinkSx: SxProps<Theme> = {
  '&:hover': {
    bgcolor: (theme) => theme.app.surfaces.cardHover
  },
  alignItems: 'center',
  borderRadius: (theme) => theme.app.radii.md,
  color: 'inherit',
  display: 'flex',
  gap: 1.5,
  justifyContent: 'space-between',
  p: 1.75,
  textDecoration: 'none',
  transition: (theme) => theme.transitions.create(['background-color', 'color'])
};

export const recentRoomsCardItemIconSx = (isActive: boolean): SxProps<Theme> => ({
  alignItems: 'center',
  bgcolor: 'background.paper',
  border: (theme) => `1px solid ${theme.app.borders.ghost}`,
  color: isActive ? 'primary.main' : 'text.secondary',
  display: 'flex',
  height: 48,
  justifyContent: 'center',
  width: 48
});

export const recentRoomsCardArrowSx: SxProps<Theme> = {
  ml: 0.5
};

export const recentRoomsCardEmptyStateSx: SxProps<Theme> = {};
