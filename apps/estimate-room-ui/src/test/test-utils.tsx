import { configureStore } from '@reduxjs/toolkit';
import { CssBaseline, ThemeProvider } from '@mui/material';
import { render, type RenderOptions } from '@testing-library/react';
import type { PropsWithChildren, ReactElement } from 'react';
import { Provider } from 'react-redux';
import { MemoryRouter, type MemoryRouterProps } from 'react-router-dom';

import { rootReducer } from '@/app/store/rootReducer';
import type { RootState } from '@/app/store/store';
import { createAppTheme } from '@/theme';

export const createTestStore = (preloadedState?: Partial<RootState>) =>
  configureStore({
    reducer: rootReducer,
    preloadedState,
    devTools: false
  });

interface ExtendedRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  readonly preloadedState?: Partial<RootState>;
  readonly routerProps?: MemoryRouterProps;
}

export const renderWithProviders = (
  ui: ReactElement,
  options: ExtendedRenderOptions = {}
) => {
  const { preloadedState, routerProps, ...renderOptions } = options;
  const store = createTestStore(preloadedState);

  const AppTestProviders = ({ children }: PropsWithChildren) => (
    <Provider store={store}>
      <MemoryRouter {...routerProps}>
        <ThemeProvider theme={createAppTheme('light')}>
          <CssBaseline />
          {children}
        </ThemeProvider>
      </MemoryRouter>
    </Provider>
  );

  return {
    store,
    ...render(ui, {
      wrapper: AppTestProviders,
      ...renderOptions
    })
  };
};

export * from '@testing-library/react';
