import type { PropsWithChildren } from 'react';

import { AuthSessionBootstrap } from '@/modules/auth';
import { AppThemeProvider } from '@/theme';

import { StoreProvider } from './StoreProvider';

export const AppProviders = ({ children }: PropsWithChildren) => (
  <StoreProvider>
    <AppThemeProvider>
      <AuthSessionBootstrap>{children}</AuthSessionBootstrap>
    </AppThemeProvider>
  </StoreProvider>
);
