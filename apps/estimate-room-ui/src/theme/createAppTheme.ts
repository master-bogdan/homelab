import { createTheme } from '@mui/material/styles';

import { getComponentOverrides, getPaletteTokens, getShadowTokens, shapeTokens, typographyTokens } from './tokens';
import type { ThemeMode } from './types';

export const createAppTheme = (mode: ThemeMode) =>
  createTheme({
    palette: getPaletteTokens(mode),
    typography: typographyTokens,
    spacing: 8,
    shape: shapeTokens,
    shadows: getShadowTokens(mode),
    components: getComponentOverrides(mode)
  });
