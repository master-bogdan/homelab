import type { PropsWithChildren } from 'react';

import { AppThemeProvider } from './AppThemeProvider';
import { StoreProvider } from './StoreProvider';

export const AppProviders = ({ children }: PropsWithChildren) => (
  <StoreProvider>
    <AppThemeProvider>{children}</AppThemeProvider>
  </StoreProvider>
);
