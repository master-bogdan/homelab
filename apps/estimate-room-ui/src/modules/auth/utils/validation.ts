import { PasswordRecommendationRules } from '../constants';

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
  const failedRule = PasswordRecommendationRules.find((rule) => !rule.test(value));

  if (failedRule) {
    return failedRule.validationMessage;
  }

  return true;
};

export const createPasswordValidationRules = (requiredMessage = 'Password is required.') => ({
  required: requiredMessage,
  validate: validatePasswordStrength
});
