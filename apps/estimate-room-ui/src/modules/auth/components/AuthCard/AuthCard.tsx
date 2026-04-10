import { Paper } from '@mui/material';
import type { PaperProps } from '@mui/material';

import { authCardRootSx } from './styles';

export type AuthCardProps = PaperProps;

export const AuthCard = ({ children, sx, ...paperProps }: AuthCardProps) => {
  const rootSx = Array.isArray(sx) ? [authCardRootSx, ...sx] : sx ? [authCardRootSx, sx] : [authCardRootSx];

  return (
    <Paper
      elevation={0}
      sx={rootSx}
      {...paperProps}
    >
      {children}
    </Paper>
  );
};
