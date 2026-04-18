import { createTheme } from '@mui/material/styles';

import {
  getAppTokens,
  getComponentOverrides,
  getPaletteTokens,
  getShadowTokens,
  getShapeTokens,
  getTypographyTokens
} from '@/shared/constants/themeTokens';
import type { ThemeMode } from '@/shared/types';

export const createAppTheme = (mode: ThemeMode) =>
  createTheme({
    app: getAppTokens(mode),
    palette: getPaletteTokens(mode),
    typography: getTypographyTokens(mode),
    spacing: 8,
    shape: getShapeTokens(mode),
    shadows: getShadowTokens(mode),
    components: getComponentOverrides(mode)
  });
