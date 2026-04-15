import type { PropsWithChildren } from 'react';
import { useMemo } from 'react';
import { CssBaseline, ThemeProvider } from '@mui/material';

import { useAppSelector } from '@/shared/store';
import { selectThemeMode } from '@/modules/system';
import { createAppTheme } from '@/shared/theme';

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
