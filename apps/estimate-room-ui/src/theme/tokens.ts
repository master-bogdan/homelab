import { createTheme } from '@mui/material/styles';
import type { Components, PaletteOptions, Shadows, Theme, ThemeOptions } from '@mui/material/styles';

import type { ThemeMode } from './types';

const brandColors = {
  info: '#0288d1',
  primary: '#1e66f5',
  secondary: '#0f9d8a',
  success: '#2e7d32',
  warning: '#ed6c02'
} as const;

const baseShadows = createTheme().shadows;

export const getPaletteTokens = (mode: ThemeMode): PaletteOptions => ({
  mode,
  primary: {
    main: brandColors.primary
  },
  secondary: {
    main: brandColors.secondary
  },
  success: {
    main: brandColors.success
  },
  info: {
    main: brandColors.info
  },
  warning: {
    main: brandColors.warning
  },
  background:
    mode === 'light'
      ? {
          default: '#f4f7fb',
          paper: '#ffffff'
        }
      : {
          default: '#0f172a',
          paper: '#111c34'
        },
  divider: mode === 'light' ? 'rgba(15, 23, 42, 0.08)' : 'rgba(255, 255, 255, 0.08)',
  text:
    mode === 'light'
      ? {
          primary: '#14213d',
          secondary: '#5c667a'
        }
      : {
          primary: '#f8fafc',
          secondary: '#c9d2e3'
        }
});

export const typographyTokens: ThemeOptions['typography'] = {
  fontFamily: '"Aptos", "Segoe UI", "Helvetica Neue", sans-serif',
  h1: {
    fontSize: '3.25rem',
    fontWeight: 700,
    lineHeight: 1.05
  },
  h2: {
    fontSize: '2.5rem',
    fontWeight: 700,
    lineHeight: 1.1
  },
  h3: {
    fontSize: '2rem',
    fontWeight: 700,
    lineHeight: 1.15
  },
  h4: {
    fontSize: '1.625rem',
    fontWeight: 700,
    lineHeight: 1.2
  },
  h5: {
    fontSize: '1.25rem',
    fontWeight: 700,
    lineHeight: 1.3
  },
  h6: {
    fontSize: '1.05rem',
    fontWeight: 700,
    lineHeight: 1.4
  },
  button: {
    fontWeight: 600,
    textTransform: 'none'
  }
};

export const shapeTokens = {
  borderRadius: 16
} as const;

export const getShadowTokens = (mode: ThemeMode): Shadows => {
  const shadows = [...baseShadows] as Shadows;

  shadows[1] =
    mode === 'light'
      ? '0 8px 20px rgba(15, 23, 42, 0.06)'
      : '0 10px 24px rgba(2, 6, 23, 0.45)';
  shadows[4] =
    mode === 'light'
      ? '0 18px 40px rgba(15, 23, 42, 0.08)'
      : '0 22px 52px rgba(2, 6, 23, 0.55)';

  return shadows;
};

export const getComponentOverrides = (mode: ThemeMode): Components<Theme> => ({
  MuiAppBar: {
    styleOverrides: {
      root: {
        backgroundImage: 'none'
      }
    }
  },
  MuiButton: {
    defaultProps: {
      disableElevation: true
    },
    styleOverrides: {
      root: ({ theme }) => ({
        borderRadius: `calc(${theme.shape.borderRadius} * 0.85)`,
        paddingInline: theme.spacing(2)
      })
    }
  },
  MuiCard: {
    styleOverrides: {
      root: {
        backgroundImage: 'none'
      }
    }
  },
  MuiCssBaseline: {
    styleOverrides: {
      body: {
        backgroundImage:
          mode === 'light'
            ? 'radial-gradient(circle at top, rgba(30, 102, 245, 0.08), transparent 40%)'
            : 'radial-gradient(circle at top, rgba(30, 102, 245, 0.12), transparent 35%)'
      }
    }
  },
  MuiPaper: {
    styleOverrides: {
      root: ({ theme }) => ({
        backgroundImage: 'none',
        borderColor: theme.palette.divider
      })
    }
  }
});
