import type { PropsWithChildren } from 'react';

import { useAuthSessionBootstrap } from '../../hooks';

export const AuthSessionBootstrap = ({ children }: PropsWithChildren) => {
  useAuthSessionBootstrap();

  return children;
};
