import type { PropsWithChildren } from 'react';

import { AppThemeProvider } from '@/theme';

import { StoreProvider } from './StoreProvider';

export const AppProviders = ({ children }: PropsWithChildren) => (
  <StoreProvider>
    <AppThemeProvider>{children}</AppThemeProvider>
  </StoreProvider>
);
