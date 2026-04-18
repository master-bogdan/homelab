import type { SxProps, Theme } from '@mui/material/styles';

export type AuthPageLayoutPattern = 'ambient' | 'dots';

export const authPageLayoutRootSx: SxProps<Theme> = {
  backgroundColor: 'background.default',
  display: 'flex',
  flexDirection: 'column',
  minHeight: '100vh',
  overflow: 'hidden',
  position: 'relative'
};

export const getAuthPageLayoutBackdropSx = (pattern: AuthPageLayoutPattern): SxProps<Theme> => ({
  '--app-auth-dot-color': (theme) => theme.app.borders.ghost,
  backgroundImage: (theme) =>
    pattern === 'dots' ? theme.app.backgrounds.authDots : theme.app.backgrounds.authAmbient,
  backgroundPosition: pattern === 'dots' ? '0 0, 0 0' : '0 0',
  backgroundSize: pattern === 'dots' ? '24px 24px, auto' : 'auto',
  inset: 0,
  opacity: pattern === 'dots' ? 0.8 : 1,
  pointerEvents: 'none',
  position: 'absolute'
});

export const authPageLayoutGlowTopSx: SxProps<Theme> = {
  bgcolor: 'primary.main',
  borderRadius: (theme) => theme.app.radii.circle,
  filter: 'blur(120px)',
  height: 320,
  opacity: 0.08,
  position: 'absolute',
  right: '-10%',
  top: '8%',
  width: 360
};

export const authPageLayoutGlowBottomSx: SxProps<Theme> = {
  bgcolor: 'secondary.main',
  borderRadius: (theme) => theme.app.radii.circle,
  bottom: '-10%',
  filter: 'blur(120px)',
  height: 260,
  left: '-8%',
  opacity: 0.22,
  position: 'absolute',
  width: 300
};

export const authPageLayoutHeaderRootSx: SxProps<Theme> = {
  px: { xs: 2.5, md: 3 },
  py: 2,
  zIndex: 1
};

export const authPageLayoutHomeLinkSx: SxProps<Theme> = {
  alignItems: 'center',
  display: 'inline-flex',
  gap: 1,
  textDecoration: 'none'
};

export const authPageLayoutMainRootSx: SxProps<Theme> = {
  alignItems: 'center',
  display: 'flex',
  flex: 1,
  justifyContent: 'center',
  px: 2,
  position: 'relative',
  py: { xs: 6, md: 8 },
  zIndex: 1
};

export const authPageLayoutInnerSx: SxProps<Theme> = {
  maxWidth: 520,
  width: '100%'
};

export const authPageLayoutFooterRootSx: SxProps<Theme> = {
  position: 'relative',
  zIndex: 1
};

export const authPageLayoutFooterStackSx: SxProps<Theme> = {
  px: { xs: 2.5, md: 3 },
  py: 2.5
};

export const authPageLayoutUtilityLinkSx: SxProps<Theme> = {
  textUnderlineOffset: 2
};
