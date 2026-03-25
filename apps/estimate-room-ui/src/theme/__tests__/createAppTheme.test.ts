import { createAppTheme } from '../createAppTheme';

describe.each([
  {
    expectedBackdropBlur: '8px',
    expectedBodyFont: 'Inter',
    expectedDisplayFont: 'Manrope',
    expectedGradient: 'linear-gradient(135deg, #5148d7 0%, #4439cb 100%)',
    expectedPrimary: '#5148d7',
    expectedSection: '#eff4ff',
    expectedTextPrimary: '#00345e',
    mode: 'light' as const
  },
  {
    expectedBackdropBlur: '20px',
    expectedBodyFont: 'Manrope',
    expectedDisplayFont: 'Manrope',
    expectedGradient: 'linear-gradient(135deg, #b4c5ff 0%, #2563eb 100%)',
    expectedPrimary: '#2563eb',
    expectedSection: '#131b2e',
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
    expectedSection,
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
      expect(theme.typography.h1?.fontFamily).toContain(expectedDisplayFont);
      expect(theme.typography.body1?.fontFamily).toContain(expectedBodyFont);
      expect(theme.shadows[12]).not.toBe('none');
    });
  }
);
