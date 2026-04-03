import type { PaperProps } from '@mui/material';
import { Paper } from '@mui/material';

export type AuthCardProps = PaperProps;

export const AuthCard = ({ children, sx, ...paperProps }: AuthCardProps) => (
  <Paper
    elevation={0}
    sx={{
      backgroundColor: 'background.paper',
      border: (theme) => `1px solid ${theme.app.borders.ghost}`,
      borderRadius: (theme) => Number(theme.shape.borderRadius) * 1.5,
      px: { xs: 3, sm: 4 },
      py: { xs: 3, sm: 4 },
      ...sx
    }}
    {...paperProps}
  >
    {children}
  </Paper>
);
