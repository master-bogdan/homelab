import { alpha, createTheme } from '@mui/material/styles';
import type { Components, PaletteOptions, Shadows, Theme, ThemeOptions } from '@mui/material/styles';

import type { AppThemeTokens, ThemeMode } from './types';

const baseShadows = createTheme().shadows;

const lightAppTokens: AppThemeTokens = {
  surfaces: {
    base: '#f8f9ff',
    section: '#eff4ff',
    card: '#ffffff',
    cardHover: '#f5f8ff',
    inset: '#e7eefc',
    overlay: alpha('#ffffff', 0.82),
    bright: '#ffffff',
    well: '#dfe9ff',
    rowAlternate: '#f3f7ff'
  },
  borders: {
    ghost: alpha('#87b5f0', 0.15),
    focusRing: alpha('#bcd0ff', 0.55)
  },
  gradients: {
    primary: 'linear-gradient(160deg, #0053db 0%, #0048c1 100%)'
  },
  effects: {
    ambientShadow: '0 12px 32px rgba(0, 52, 94, 0.06)',
    backdropBlur: '12px'
  },
  layout: {
    pageGap: 16,
    sectionGap: 12
  }
};

const darkAppTokens: AppThemeTokens = {
  surfaces: {
    base: '#0b1326',
    section: '#131b2e',
    card: '#171f33',
    cardHover: '#1d2640',
    inset: '#0f182a',
    overlay: alpha('#31394d', 0.6),
    bright: '#31394d',
    well: '#2d3449',
    rowAlternate: '#12192c'
  },
  borders: {
    ghost: alpha('#434655', 0.15),
    focusRing: alpha('#2563eb', 0.28)
  },
  gradients: {
    primary: 'linear-gradient(135deg, #b4c5ff 0%, #2563eb 100%)'
  },
  effects: {
    ambientShadow: '0 12px 32px rgba(0, 0, 0, 0.25)',
    backdropBlur: '20px'
  },
  layout: {
    pageGap: 16,
    sectionGap: 12
  }
};

export const getAppTokens = (mode: ThemeMode): AppThemeTokens =>
  mode === 'light' ? lightAppTokens : darkAppTokens;

export const getPaletteTokens = (mode: ThemeMode): PaletteOptions => {
  const appTokens = getAppTokens(mode);

  if (mode === 'light') {
    return {
      mode,
      primary: {
        contrastText: '#f7f9ff',
        dark: '#0048c1',
        light: '#4d84eb',
        main: '#0053db'
      },
      secondary: {
        contrastText: '#173057',
        dark: '#d5e1ff',
        light: '#f4f7ff',
        main: '#e3ebff'
      },
      success: {
        contrastText: '#f7fffb',
        dark: '#0e6b54',
        light: '#4aa98b',
        main: '#17886b'
      },
      info: {
        contrastText: '#f5f9ff',
        dark: '#2054a8',
        light: '#77a3f0',
        main: '#2f73da'
      },
      warning: {
        contrastText: '#fffaf1',
        dark: '#a56913',
        light: '#d8ab58',
        main: '#c18624'
      },
      error: {
        contrastText: '#fff7f6',
        dark: '#a54843',
        light: '#dd8f89',
        main: '#c4615a'
      },
      background: {
        default: appTokens.surfaces.base,
        paper: appTokens.surfaces.card
      },
      divider: appTokens.borders.ghost,
      action: {
        active: '#224271',
        focus: alpha('#0053db', 0.12),
        hover: '#e6eeff',
        selected: alpha('#0053db', 0.1)
      },
      text: {
        primary: '#102440',
        secondary: '#556b88'
      }
    };
  }

  return {
    mode,
    primary: {
      contrastText: '#eeefff',
      dark: '#1e4fc0',
      light: '#b4c5ff',
      main: '#2563eb'
    },
    secondary: {
      contrastText: '#a4b6f5',
      dark: '#2c3c6f',
      light: '#3c5293',
      main: '#33467e'
    },
    success: {
      contrastText: '#072012',
      dark: '#4fbc7f',
      light: '#a7edc3',
      main: '#72d69c'
    },
    info: {
      contrastText: '#071524',
      dark: '#6aa8ff',
      light: '#cae0ff',
      main: '#8fc0ff'
    },
    warning: {
      contrastText: '#211100',
      dark: '#efb85d',
      light: '#ffe0a6',
      main: '#f4ca84'
    },
    error: {
      contrastText: '#240807',
      dark: '#ff8c81',
      light: '#ffd4cf',
      main: '#ffb4ab'
    },
    background: {
      default: appTokens.surfaces.base,
      paper: appTokens.surfaces.card
    },
    divider: appTokens.borders.ghost,
    action: {
      active: '#dae2fd',
      focus: alpha('#2563eb', 0.24),
      hover: '#31394d',
      selected: alpha('#b4c5ff', 0.16)
    },
    text: {
      primary: '#dae2fd',
      secondary: '#c3c6d7'
    }
  };
};

export const getTypographyTokens = (mode: ThemeMode): ThemeOptions['typography'] => {
  const displayFont = '"Manrope", "Inter", "Segoe UI", sans-serif';
  const bodyFont =
    mode === 'light'
      ? '"Inter", "Manrope", "Segoe UI", sans-serif'
      : '"Manrope", "Inter", "Segoe UI", sans-serif';

  return {
    fontFamily: bodyFont,
    h1: {
      fontFamily: displayFont,
      fontSize: '3.5rem',
      fontWeight: 800,
      letterSpacing: '-0.02em',
      lineHeight: 1.02
    },
    h2: {
      fontFamily: displayFont,
      fontSize: '2.75rem',
      fontWeight: 800,
      letterSpacing: '-0.02em',
      lineHeight: 1.06
    },
    h3: {
      fontFamily: displayFont,
      fontSize: '2rem',
      fontWeight: 750,
      letterSpacing: '-0.02em',
      lineHeight: 1.12
    },
    h4: {
      fontFamily: displayFont,
      fontSize: '1.5rem',
      fontWeight: 750,
      letterSpacing: '-0.02em',
      lineHeight: 1.18
    },
    h5: {
      fontFamily: bodyFont,
      fontSize: '1.125rem',
      fontWeight: 700,
      lineHeight: 1.28
    },
    h6: {
      fontFamily: bodyFont,
      fontSize: '1rem',
      fontWeight: 700,
      lineHeight: 1.36
    },
    subtitle1: {
      fontFamily: bodyFont,
      fontSize: '1rem',
      fontWeight: 600,
      lineHeight: 1.5
    },
    subtitle2: {
      fontFamily: bodyFont,
      fontSize: '0.875rem',
      fontWeight: 700,
      lineHeight: 1.45
    },
    body1: {
      fontFamily: bodyFont,
      fontSize: '1rem',
      lineHeight: 1.65
    },
    body2: {
      fontFamily: bodyFont,
      fontSize: '0.9375rem',
      lineHeight: 1.6
    },
    caption: {
      fontFamily: bodyFont,
      fontSize: '0.8125rem',
      lineHeight: 1.45
    },
    overline: {
      fontFamily: bodyFont,
      fontSize: '0.6875rem',
      fontWeight: 700,
      letterSpacing: '0.05em',
      lineHeight: 1.4,
      textTransform: 'uppercase'
    },
    button: {
      fontFamily: bodyFont,
      fontWeight: 700,
      letterSpacing: '0.01em',
      textTransform: 'none'
    }
  };
};

export const getShapeTokens = (mode: ThemeMode) =>
  ({
    borderRadius: mode === 'light' ? 6 : 4
  }) as const;

export const getShadowTokens = (mode: ThemeMode): Shadows => {
  const shadows = [...baseShadows] as Shadows;
  const ambientShadow = getAppTokens(mode).effects.ambientShadow;

  shadows.fill('none');
  shadows[8] = ambientShadow;
  shadows[12] = ambientShadow;
  shadows[16] = ambientShadow;
  shadows[24] = ambientShadow;

  return shadows;
};

export const getComponentOverrides = (mode: ThemeMode): Components<Theme> => ({
  MuiAppBar: {
    styleOverrides: {
      root: ({ theme }) => ({
        backdropFilter: `blur(${theme.app.effects.backdropBlur})`,
        backgroundColor: theme.app.surfaces.overlay,
        backgroundImage: 'none',
        boxShadow: 'none'
      })
    }
  },
  MuiAlert: {
    styleOverrides: {
      root: ({ theme }) => ({
        alignItems: 'center',
        backgroundColor: theme.app.surfaces.section,
        border: 'none',
        borderRadius: theme.shape.borderRadius * 2,
        boxShadow: 'none'
      }),
      standardError: ({ theme }) => ({
        borderLeft: `4px solid ${theme.palette.error.main}`,
        color: theme.palette.text.primary
      }),
      standardInfo: ({ theme }) => ({
        borderLeft: `4px solid ${theme.palette.info.main}`,
        color: theme.palette.text.primary
      }),
      standardSuccess: ({ theme }) => ({
        borderLeft: `4px solid ${theme.palette.success.main}`,
        color: theme.palette.text.primary
      }),
      standardWarning: ({ theme }) => ({
        borderLeft: `4px solid ${theme.palette.warning.main}`,
        color: theme.palette.text.primary
      })
    }
  },
  MuiButton: {
    defaultProps: {
      disableElevation: true
    },
    styleOverrides: {
      root: ({ theme }) => ({
        borderRadius: theme.shape.borderRadius,
        minHeight: 40,
        paddingInline: theme.spacing(2),
        transition: theme.transitions.create(['background-color', 'border-color', 'color'])
      }),
      contained: ({ theme }) => ({
        boxShadow: 'none',
        '&:hover': {
          boxShadow: 'none'
        },
        '&.Mui-disabled': {
          backgroundImage: 'none'
        }
      }),
      containedPrimary: ({ theme }) => ({
        backgroundColor: 'transparent',
        backgroundImage: theme.app.gradients.primary,
        color: theme.palette.primary.contrastText,
        '&:hover': {
          backgroundColor: 'transparent',
          backgroundImage: theme.app.gradients.primary,
          filter: 'brightness(1.04)'
        }
      }),
      containedSecondary: ({ theme }) => ({
        backgroundColor: theme.palette.secondary.main,
        color: theme.palette.secondary.contrastText,
        '&:hover': {
          backgroundColor: theme.palette.secondary.dark
        }
      }),
      outlined: ({ theme }) => ({
        backgroundColor: 'transparent',
        borderColor: theme.app.borders.ghost,
        color: theme.palette.text.primary,
        '&:hover': {
          backgroundColor: theme.app.surfaces.section,
          borderColor: theme.app.borders.ghost
        }
      }),
      text: ({ theme }) => ({
        color:
          theme.palette.mode === 'light'
            ? theme.palette.text.secondary
            : theme.palette.primary.light,
        '&:hover': {
          backgroundColor: theme.app.surfaces.section
        }
      })
    }
  },
  MuiCard: {
    styleOverrides: {
      root: ({ theme }) => ({
        backgroundColor: theme.app.surfaces.card,
        backgroundImage: 'none',
        border: 'none',
        boxShadow: 'none'
      })
    }
  },
  MuiChip: {
    styleOverrides: {
      root: ({ theme }) => ({
        borderRadius: theme.shape.borderRadius,
        fontWeight: 600
      }),
      filled: ({ theme }) => ({
        backgroundColor: theme.app.surfaces.section,
        color: theme.palette.text.primary
      }),
      filledPrimary: ({ theme }) => ({
        backgroundColor:
          theme.palette.mode === 'light'
            ? alpha(theme.palette.primary.main, 0.12)
            : theme.palette.secondary.main,
        color:
          theme.palette.mode === 'light'
            ? theme.palette.primary.main
            : theme.palette.secondary.contrastText
      }),
      outlined: ({ theme }) => ({
        backgroundColor: 'transparent',
        borderColor: theme.app.borders.ghost,
        color: theme.palette.text.secondary
      })
    }
  },
  MuiCssBaseline: {
    styleOverrides: {
      body: {
        backgroundImage:
          mode === 'light'
            ? 'radial-gradient(circle at top left, rgba(0, 83, 219, 0.12), transparent 32%), linear-gradient(180deg, #f8f9ff 0%, #eff4ff 100%)'
            : 'radial-gradient(circle at top left, rgba(37, 99, 235, 0.24), transparent 34%), linear-gradient(180deg, #0b1326 0%, #10182d 100%)'
      },
      '::selection': {
        backgroundColor:
          mode === 'light' ? alpha('#0053db', 0.16) : alpha('#b4c5ff', 0.24)
      }
    }
  },
  MuiDivider: {
    styleOverrides: {
      root: ({ theme }) => ({
        borderColor: theme.app.borders.ghost
      })
    }
  },
  MuiDrawer: {
    styleOverrides: {
      paper: ({ theme }) => ({
        backgroundColor: theme.app.surfaces.section,
        backgroundImage:
          theme.palette.mode === 'light'
            ? 'linear-gradient(180deg, rgba(255, 255, 255, 0.45) 0%, rgba(239, 244, 255, 0.85) 100%)'
            : 'linear-gradient(180deg, rgba(49, 57, 77, 0.35) 0%, rgba(19, 27, 46, 1) 100%)',
        border: 'none'
      })
    }
  },
  MuiInputBase: {
    styleOverrides: {
      input: ({ theme }) => ({
        paddingBlock: theme.spacing(1.5)
      })
    }
  },
  MuiListItemButton: {
    styleOverrides: {
      root: ({ theme }) => ({
        borderRadius: theme.shape.borderRadius * 2,
        marginBottom: theme.spacing(0.5),
        paddingBlock: theme.spacing(1.1),
        '&.Mui-selected': {
          backgroundColor: theme.app.surfaces.card,
          color: theme.palette.text.primary,
          '&:hover': {
            backgroundColor: theme.app.surfaces.card
          }
        },
        '&:hover': {
          backgroundColor: theme.app.surfaces.rowAlternate
        }
      })
    }
  },
  MuiOutlinedInput: {
    styleOverrides: {
      notchedOutline: ({ theme }) => ({
        borderColor: theme.app.borders.ghost
      }),
      root: ({ theme }) => ({
        backgroundColor:
          theme.palette.mode === 'light'
            ? theme.app.surfaces.card
            : theme.app.surfaces.section,
        borderRadius: theme.shape.borderRadius,
        transition: theme.transitions.create(['background-color', 'box-shadow', 'border-color']),
        '&:hover .MuiOutlinedInput-notchedOutline': {
          borderColor: alpha(theme.palette.primary.main, 0.3)
        },
        '&.Mui-focused': {
          backgroundColor: theme.app.surfaces.card,
          boxShadow: `0 0 0 2px ${theme.app.borders.focusRing}`
        },
        '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
          borderColor: alpha(theme.palette.primary.main, 0.5)
        },
        '&.Mui-error .MuiOutlinedInput-notchedOutline': {
          borderColor: alpha(theme.palette.error.main, 0.7)
        }
      })
    }
  },
  MuiPaper: {
    styleOverrides: {
      root: ({ theme }) => ({
        backgroundColor: theme.palette.background.paper,
        backgroundImage: 'none',
        border: 'none',
        boxShadow: 'none'
      })
    }
  },
  MuiTooltip: {
    styleOverrides: {
      tooltip: ({ theme }) => ({
        backdropFilter: `blur(${theme.app.effects.backdropBlur})`,
        backgroundColor: theme.app.surfaces.overlay,
        boxShadow: theme.app.effects.ambientShadow,
        color: theme.palette.text.primary
      })
    }
  }
});
