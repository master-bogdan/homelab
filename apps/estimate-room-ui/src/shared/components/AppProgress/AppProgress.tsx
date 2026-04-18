import { CircularProgress, LinearProgress } from '@mui/material';
import type { CircularProgressProps, LinearProgressProps } from '@mui/material';

export type AppCircularProgressProps = CircularProgressProps & {
  readonly kind?: 'circular';
};

export type AppLinearProgressProps = LinearProgressProps & {
  readonly kind: 'linear';
};

export type AppProgressProps = AppCircularProgressProps | AppLinearProgressProps;

export const AppProgress = ({ kind = 'circular', ...progressProps }: AppProgressProps) =>
  kind === 'linear' ? (
    <LinearProgress {...(progressProps as LinearProgressProps)} />
  ) : (
    <CircularProgress {...(progressProps as CircularProgressProps)} />
  );
