import type { PropsWithChildren } from 'react';
import { useMemo } from 'react';
import { CssBaseline, ThemeProvider } from '@mui/material';

import { useAppSelector } from '@/app/store/hooks';
import { selectThemeMode } from '@/app/store/uiSelectors';

import { createAppTheme } from './createAppTheme';

export const AppThemeProvider = ({ children }: PropsWithChildren) => {
  const themeMode = useAppSelector(selectThemeMode);
  const theme = useMemo(() => createAppTheme(themeMode), [themeMode]);

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      {children}
    </ThemeProvider>
  );
};
