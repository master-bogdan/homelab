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
  backgroundImage:
    pattern === 'dots'
      ? (theme) =>
          `radial-gradient(${theme.app.borders.ghost} 0.6px, transparent 0.6px), linear-gradient(90deg, rgba(81, 72, 215, 0.04) 0%, transparent 36%, rgba(81, 72, 215, 0.05) 72%, transparent 100%)`
      : 'linear-gradient(90deg, rgba(81, 72, 215, 0.04) 0%, transparent 36%, rgba(81, 72, 215, 0.05) 72%, transparent 100%)',
  backgroundPosition: pattern === 'dots' ? '0 0, 0 0' : '0 0',
  backgroundSize: pattern === 'dots' ? '24px 24px, auto' : 'auto',
  inset: 0,
  opacity: pattern === 'dots' ? 0.8 : 1,
  pointerEvents: 'none',
  position: 'absolute'
});

export const authShellGlowTopSx: SxProps<Theme> = {
  bgcolor: 'primary.main',
  borderRadius: '50%',
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
  borderRadius: '50%',
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
