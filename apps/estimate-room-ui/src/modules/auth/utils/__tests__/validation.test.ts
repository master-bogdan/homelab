import {
  createEmailValidationRules,
  normalizeEmailAddress,
  validateEmailAddress,
  validatePasswordStrength
} from '../validation';

describe('auth validation utilities', () => {
  it('trims email input before validation', () => {
    const rules = createEmailValidationRules();

    expect(normalizeEmailAddress('  name@company.com  ')).toBe('name@company.com');
    expect(rules.setValueAs('  name@company.com  ')).toBe('name@company.com');
  });

  it('rejects malformed email addresses', () => {
    expect(validateEmailAddress('invalid-address')).toBe('Enter a valid email address.');
    expect(validateEmailAddress('name@company.com')).toBe(true);
  });

  it('enforces the register and reset password requirements', () => {
    expect(validatePasswordStrength('Short1!')).toBe('Password must be at least 8 characters.');
    expect(validatePasswordStrength('Password!')).toBe('Password must include at least one number.');
    expect(validatePasswordStrength('password1!')).toBe(
      'Password must include at least one uppercase letter.'
    );
    expect(validatePasswordStrength('Password1')).toBe(
      'Password must include at least one special character.'
    );
    expect(validatePasswordStrength('Password1!')).toBe(true);
  });
});
