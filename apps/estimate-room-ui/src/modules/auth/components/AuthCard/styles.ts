import { styled } from '@mui/material';

import { AppSurface } from '@/shared/ui';

export const AuthCardRoot = styled(AppSurface)(({ theme }) => ({
  backgroundColor: theme.palette.background.paper,
  border: `1px solid ${theme.app.borders.ghost}`,
  borderRadius: theme.app.radii.lg,
  paddingBlock: theme.spacing(3),
  paddingInline: theme.spacing(3),
  [theme.breakpoints.up('sm')]: {
    paddingBlock: theme.spacing(4),
    paddingInline: theme.spacing(4)
  }
}));
