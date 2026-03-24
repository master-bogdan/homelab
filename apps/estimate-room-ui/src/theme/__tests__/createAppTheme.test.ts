import { createAppTheme } from '../createAppTheme';

describe.each([
  {
    expectedBodyFont: 'Inter',
    expectedDisplayFont: 'Manrope',
    expectedGradient: 'linear-gradient(160deg, #0053db 0%, #0048c1 100%)',
    expectedPrimary: '#0053db',
    expectedSection: '#eff4ff',
    mode: 'light' as const
  },
  {
    expectedBodyFont: 'Manrope',
    expectedDisplayFont: 'Manrope',
    expectedGradient: 'linear-gradient(135deg, #b4c5ff 0%, #2563eb 100%)',
    expectedPrimary: '#2563eb',
    expectedSection: '#131b2e',
    mode: 'dark' as const
  }
])('createAppTheme($mode)', ({ expectedBodyFont, expectedDisplayFont, expectedGradient, expectedPrimary, expectedSection, mode }) => {
  it('maps the design tokens into the MUI theme contract', () => {
    const theme = createAppTheme(mode);

    expect(theme.palette.primary.main).toBe(expectedPrimary);
    expect(theme.app.surfaces.section).toBe(expectedSection);
    expect(theme.app.gradients.primary).toBe(expectedGradient);
    expect(theme.typography.h1?.fontFamily).toContain(expectedDisplayFont);
    expect(theme.typography.body1?.fontFamily).toContain(expectedBodyFont);
    expect(theme.shadows[12]).not.toBe('none');
  });
});
