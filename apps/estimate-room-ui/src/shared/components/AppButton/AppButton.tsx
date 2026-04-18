import { Button } from '@mui/material';
import type { ButtonProps } from '@mui/material';
import type { ReactNode } from 'react';
import type { LinkProps as RouterLinkProps } from 'react-router-dom';

import { AppProgress } from '../AppProgress';

export interface AppButtonProps extends ButtonProps {
  readonly loading?: boolean;
  readonly loadingIndicator?: ReactNode;
  readonly loadingText?: ReactNode;
  readonly to?: RouterLinkProps['to'];
}

export const AppButton = ({
  children,
  disabled,
  endIcon,
  loading = false,
  loadingIndicator,
  loadingText,
  startIcon,
  ...buttonProps
}: AppButtonProps) => (
  <Button
    aria-busy={loading || undefined}
    disabled={disabled || loading}
    endIcon={loading ? undefined : endIcon}
    startIcon={
      loading
        ? (loadingIndicator ?? <AppProgress color="inherit" size={16} />)
        : startIcon
    }
    {...buttonProps}
  >
    {loading && loadingText ? loadingText : children}
  </Button>
);
