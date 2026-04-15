import type { SxProps, Theme } from '@mui/material/styles';

export type AuthShellPattern = 'ambient' | 'dots';

export const authShellRootSx: SxProps<Theme> = {
  backgroundColor: 'background.default',
  display: 'flex',
  flexDirection: 'column',
  minHeight: '100vh',
  overflow: 'hidden',
  position: 'relative'
};

export const getAuthShellBackdropSx = (pattern: AuthShellPattern): SxProps<Theme> => ({
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

export const authShellGlowTopSx: SxProps<Theme> = {
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

export const authShellGlowBottomSx: SxProps<Theme> = {
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

export const authShellHeaderRootSx: SxProps<Theme> = {
  px: { xs: 2.5, md: 3 },
  py: 2,
  zIndex: 1
};

export const authShellHomeLinkSx: SxProps<Theme> = {
  alignItems: 'center',
  display: 'inline-flex',
  gap: 1,
  textDecoration: 'none'
};

export const authShellMainRootSx: SxProps<Theme> = {
  alignItems: 'center',
  display: 'flex',
  flex: 1,
  justifyContent: 'center',
  px: 2,
  position: 'relative',
  py: { xs: 6, md: 8 },
  zIndex: 1
};

export const authShellInnerSx: SxProps<Theme> = {
  maxWidth: 520,
  width: '100%'
};

export const authShellFooterRootSx: SxProps<Theme> = {
  position: 'relative',
  zIndex: 1
};

export const authShellFooterStackSx: SxProps<Theme> = {
  px: { xs: 2.5, md: 3 },
  py: 2.5
};

export const authShellUtilityLinkSx: SxProps<Theme> = {
  textUnderlineOffset: 2
};
