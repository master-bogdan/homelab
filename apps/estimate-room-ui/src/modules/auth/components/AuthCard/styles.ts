import type { SxProps, Theme } from '@mui/material/styles';

export const authCardRootSx: SxProps<Theme> = {
  backgroundColor: 'background.paper',
  border: (theme) => `1px solid ${theme.app.borders.ghost}`,
  borderRadius: (theme) => theme.app.radii.lg,
  px: { xs: 3, sm: 4 },
  py: { xs: 3, sm: 4 }
};
