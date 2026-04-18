import type { AppSurfaceProps } from '@/shared/components';

import { AuthCardRoot } from './styles';

type AuthCardProps = AppSurfaceProps;

export const AuthCard = ({ children, ...paperProps }: AuthCardProps) => (
  <AuthCardRoot {...paperProps}>{children}</AuthCardRoot>
);
