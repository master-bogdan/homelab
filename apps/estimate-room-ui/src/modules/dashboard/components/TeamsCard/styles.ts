import type { SxProps, Theme } from '@mui/material/styles';

export const teamsCardRootSx = (isEmptyOrError: boolean): SxProps<Theme> => ({
  bgcolor: (theme) => theme.app.surfaces.section,
  borderRadius: 3,
  minHeight: 260,
  p: isEmptyOrError ? 3 : 1
});

export const teamsCardItemLinkSx: SxProps<Theme> = {
  '&:hover': {
    bgcolor: (theme) => theme.app.surfaces.cardHover
  },
  alignItems: 'center',
  borderRadius: 2,
  color: 'inherit',
  display: 'flex',
  gap: 1.5,
  justifyContent: 'space-between',
  p: 1.5,
  textDecoration: 'none'
};

export const teamsCardAvatarSx = (index: number): SxProps<Theme> => ({
  alignItems: 'center',
  bgcolor:
    index % 3 === 0
      ? 'primary.light'
      : index % 3 === 1
        ? 'secondary.main'
        : 'warning.light',
  color: index % 3 === 1 ? 'secondary.contrastText' : 'text.primary',
  display: 'flex',
  fontSize: '0.875rem',
  fontWeight: 700,
  height: 32,
  justifyContent: 'center',
  width: 32
});
