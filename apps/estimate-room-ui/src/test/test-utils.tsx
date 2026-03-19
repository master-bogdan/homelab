import { CssBaseline, ThemeProvider } from '@mui/material';
import { render, type RenderOptions } from '@testing-library/react';
import type { PropsWithChildren, ReactElement } from 'react';
import { MemoryRouter } from 'react-router-dom';

import { createAppTheme } from '@/theme';

const AppTestProviders = ({ children }: PropsWithChildren) => (
  <MemoryRouter>
    <ThemeProvider theme={createAppTheme('light')}>
      <CssBaseline />
      {children}
    </ThemeProvider>
  </MemoryRouter>
);

export const renderWithProviders = (
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>
) =>
  render(ui, {
    wrapper: AppTestProviders,
    ...options
  });

export * from '@testing-library/react';
