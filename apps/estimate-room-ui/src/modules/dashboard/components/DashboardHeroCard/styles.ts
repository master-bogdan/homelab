import type { SxProps, Theme } from '@mui/material/styles';

export const dashboardHeroCardEmptyRootSx: SxProps<Theme> = {
  border: (theme) => `1px solid ${theme.app.borders.ghost}`,
  display: 'flex',
  flexDirection: 'column',
  gap: 3,
  height: '100%',
  justifyContent: 'center',
  minHeight: 320,
  p: { xs: 3, md: 4 },
  textAlign: 'center'
};

export const dashboardHeroCardEmptyAvatarSx: SxProps<Theme> = {
  bgcolor: 'secondary.light',
  color: 'text.secondary',
  height: 64,
  width: 64
};

export const dashboardHeroCardRootSx: SxProps<Theme> = {
  border: (theme) => `1px solid ${theme.app.borders.ghost}`,
  display: 'flex',
  flexDirection: 'column',
  gap: 3,
  height: '100%',
  minHeight: 320,
  p: { xs: 3, md: 4 }
};

export const dashboardHeroCardParticipantAvatarSx = (index: number): SxProps<Theme> => ({
  bgcolor: index % 2 === 0 ? 'primary.main' : 'secondary.main',
  border: '2px solid',
  borderColor: 'background.paper',
  color: index % 2 === 0 ? 'primary.contrastText' : 'secondary.contrastText',
  ml: index === 0 ? 0 : -1
});

export const dashboardHeroCardMemberOverflowSx: SxProps<Theme> = {
  bgcolor: 'secondary.light',
  color: 'text.primary',
  ml: -1
};

export const dashboardHeroCardMetricGridSx: SxProps<Theme> = {
  display: 'grid',
  gap: 1.5,
  gridTemplateColumns: {
    xs: 'repeat(1, minmax(0, 1fr))',
    sm: 'repeat(3, minmax(0, 1fr))'
  }
};

export const dashboardHeroCardMetricCardSx: SxProps<Theme> = {
  bgcolor: 'secondary.light',
  borderRadius: 2,
  p: 1.75
};

export const dashboardHeroCardProgressBarSx: SxProps<Theme> = {
  borderRadius: 999,
  height: 8
};
