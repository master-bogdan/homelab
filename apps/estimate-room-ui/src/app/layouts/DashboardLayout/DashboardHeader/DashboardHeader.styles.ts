import type { SxProps, Theme } from '@mui/material';

export const dashboardHeaderRootSx: SxProps<Theme> = {
  backdropFilter: (theme) => `blur(${theme.app.effects.backdropBlur})`,
  bgcolor: (theme) => theme.app.surfaces.overlay,
  borderBottom: (theme) => `1px solid ${theme.app.borders.ghost}`,
  position: 'sticky',
  top: 0,
  zIndex: (theme) => theme.zIndex.appBar
};

export const dashboardHeaderToolbarSx: SxProps<Theme> = {
  alignItems: { xs: 'flex-start', md: 'center' },
  flexWrap: 'wrap',
  gap: 2,
  minHeight: '72px !important',
  px: { xs: 2, md: 4 },
  py: 1.5
};

export const dashboardHeaderMenuButtonSx: SxProps<Theme> = {
  mt: { xs: 0.5, md: 0 }
};

export const dashboardHeaderTitleSx: SxProps<Theme> = {
  flex: 1
};

export const dashboardHeaderActionsSx: SxProps<Theme> = {
  ml: 'auto'
};

export const dashboardHeaderDividerSx: SxProps<Theme> = {
  alignSelf: 'stretch',
  bgcolor: (theme) => theme.app.borders.ghost,
  display: { xs: 'none', sm: 'block' },
  width: '1px'
};

export const dashboardHeaderAvatarButtonSx: SxProps<Theme> = {
  p: 0
};

export const dashboardHeaderAvatarSx: SxProps<Theme> = {
  bgcolor: 'primary.main',
  color: 'primary.contrastText',
  height: 36,
  width: 36
};
