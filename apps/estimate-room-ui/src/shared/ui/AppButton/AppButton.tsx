import { Button, CircularProgress } from '@mui/material';
import type { ButtonProps } from '@mui/material';
import type { ReactNode } from 'react';
import type { LinkProps as RouterLinkProps } from 'react-router-dom';

export interface AppButtonProps extends ButtonProps {
  readonly loading?: boolean;
  readonly loadingIndicator?: ReactNode;
  readonly loadingText?: ReactNode;
  readonly to?: RouterLinkProps['to'];
}

export const AppButton = ({
  children,
  disabled,
  loading = false,
  loadingIndicator,
  loadingText,
  startIcon,
  ...buttonProps
}: AppButtonProps) => (
  <Button
    aria-busy={loading || undefined}
    disabled={disabled || loading}
    startIcon={
      loading
        ? (loadingIndicator ?? <CircularProgress color="inherit" size={16} />)
        : startIcon
    }
    {...buttonProps}
  >
    {loading && loadingText ? loadingText : children}
  </Button>
);
