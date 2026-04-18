import { createAppTheme } from '../createAppTheme';

describe.each([
  {
    expectedBackdropBlur: '8px',
    expectedBodyFont: 'Inter',
    expectedDisplayFont: 'Manrope',
    expectedGradient: 'linear-gradient(135deg, #5148d7 0%, #4439cb 100%)',
    expectedPrimary: '#5148d7',
    expectedRadius: 12,
    expectedSection: '#eff4ff',
    expectedStateLayer: '#eff4ff',
    expectedTextPrimary: '#00345e',
    mode: 'light' as const
  },
  {
    expectedBackdropBlur: '20px',
    expectedBodyFont: 'Manrope',
    expectedDisplayFont: 'Manrope',
    expectedGradient: 'linear-gradient(135deg, #b4c5ff 0%, #2563eb 100%)',
    expectedPrimary: '#2563eb',
    expectedRadius: 8,
    expectedSection: '#131b2e',
    expectedStateLayer: '#3c5293',
    expectedTextPrimary: '#dae2fd',
    mode: 'dark' as const
  }
])(
  'createAppTheme($mode)',
  ({
    expectedBackdropBlur,
    expectedBodyFont,
    expectedDisplayFont,
    expectedGradient,
    expectedPrimary,
    expectedRadius,
    expectedSection,
    expectedStateLayer,
    expectedTextPrimary,
    mode
  }) => {
    it('maps the design tokens into the MUI theme contract', () => {
      const theme = createAppTheme(mode);

      expect(theme.palette.primary.main).toBe(expectedPrimary);
      expect(theme.palette.text.primary).toBe(expectedTextPrimary);
      expect(theme.app.surfaces.section).toBe(expectedSection);
      expect(theme.app.gradients.primary).toBe(expectedGradient);
      expect(theme.app.effects.backdropBlur).toBe(expectedBackdropBlur);
      expect(theme.app.backgrounds.body).toContain('radial-gradient');
      expect(theme.app.backgrounds.drawer).toContain('linear-gradient');
      expect(theme.app.radii.lg).toBe(expectedRadius);
      expect(theme.app.stateLayers.secondaryPanel).toBe(expectedStateLayer);
      expect(theme.app.layout.drawerWidth).toBe(288);
      expect(theme.typography.h1?.fontFamily).toContain(expectedDisplayFont);
      expect(theme.typography.body1?.fontFamily).toContain(expectedBodyFont);
      expect(theme.shadows[12]).not.toBe('none');
    });

    it('uses transition-based button hovers without hover lift effects', () => {
      const theme = createAppTheme(mode);
      const buttonOverrides = theme.components?.MuiButton?.styleOverrides;
      const rootStyles = (buttonOverrides?.root as (args: { theme: typeof theme }) => Record<string, string>)({
        theme
      });
      const containedPrimaryStyles = (
        buttonOverrides?.containedPrimary as (args: { theme: typeof theme }) => Record<string, unknown>
      )({ theme });
      const hoverStyles = (
        containedPrimaryStyles['@media (hover: hover)'] as Record<string, Record<string, unknown>>
      )['&:hover'];

      expect(rootStyles.transition).toContain('opacity');
      expect(rootStyles.transition).not.toContain('filter');
      expect(rootStyles.transition).not.toContain('transform');
      expect(containedPrimaryStyles.backgroundBlendMode).toBe('soft-light');
      expect(hoverStyles.filter).toBeUndefined();
      expect(hoverStyles.transform).toBeUndefined();
    });

    it('applies shared accessibility defaults for focus and helper text', () => {
      const theme = createAppTheme(mode);
      const iconButtonOverrides = theme.components?.MuiIconButton?.styleOverrides;
      const helperTextOverrides = theme.components?.MuiFormHelperText?.styleOverrides;
      const iconButtonStyles = (
        iconButtonOverrides?.root as (args: { theme: typeof theme }) => Record<string, unknown>
      )({ theme });
      const helperTextStyles = (
        helperTextOverrides?.root as (args: { theme: typeof theme }) => Record<string, unknown>
      )({ theme });

      expect(iconButtonStyles['&:focus-visible']).toBeDefined();
      expect(iconButtonStyles['&.Mui-disabled']).toBeDefined();
      expect(helperTextStyles.minHeight).toBe(theme.spacing(2.5));
    });
  }
);
