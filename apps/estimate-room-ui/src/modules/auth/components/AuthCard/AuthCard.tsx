import type { AppSurfaceProps } from '@/shared/ui';

import { AuthCardRoot } from './styles';

type AuthCardProps = AppSurfaceProps;

export const AuthCard = ({ children, ...paperProps }: AuthCardProps) => (
  <AuthCardRoot {...paperProps}>{children}</AuthCardRoot>
);
