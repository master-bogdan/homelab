import { Paper } from '@mui/material';
import type { PaperProps } from '@mui/material';

export type AppSurfaceProps = PaperProps;

export const AppSurface = ({ elevation = 0, ...paperProps }: AppSurfaceProps) => (
  <Paper elevation={elevation} {...paperProps} />
);
