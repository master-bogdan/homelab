const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/u;

export const normalizeEmailAddress = (value: string) => value.trim();

export const validateEmailAddress = (value: string) =>
  emailPattern.test(value) || 'Enter a valid email address.';

export const createEmailValidationRules = (requiredMessage = 'Email is required.') => ({
  required: requiredMessage,
  setValueAs: (value: string) =>
    typeof value === 'string' ? normalizeEmailAddress(value) : value,
  validate: validateEmailAddress
});

export const validatePasswordStrength = (value: string) => {
  if (value.length < 8) {
    return 'Password must be at least 8 characters.';
  }

  if (!/[0-9]/u.test(value)) {
    return 'Password must include at least one number.';
  }

  if (!/[A-Z]/u.test(value)) {
    return 'Password must include at least one uppercase letter.';
  }

  if (!/[^A-Za-z0-9]/u.test(value)) {
    return 'Password must include at least one special character.';
  }

  return true;
};

export const createPasswordValidationRules = (requiredMessage = 'Password is required.') => ({
  required: requiredMessage,
  validate: validatePasswordStrength
});
