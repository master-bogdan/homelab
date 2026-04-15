import { AppSurface } from '@/shared/ui';
import type { AppSurfaceProps } from '@/shared/ui';

import { authCardRootSx } from './styles';

export type AuthCardProps = AppSurfaceProps;

export const AuthCard = ({ children, sx, ...paperProps }: AuthCardProps) => {
  const rootSx = Array.isArray(sx) ? [authCardRootSx, ...sx] : sx ? [authCardRootSx, sx] : [authCardRootSx];

  return (
    <AppSurface
      sx={rootSx}
      {...paperProps}
    >
      {children}
    </AppSurface>
  );
};
