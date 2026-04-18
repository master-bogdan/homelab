export const PasswordRecommendationRuleIds = {
  MIN_LENGTH: 'minLength',
  NUMBER: 'number',
  SPECIAL_CHARACTER: 'specialCharacter',
  UPPERCASE: 'uppercase'
} as const;

export const PasswordMinLength = 8;

export const PasswordRecommendationRules = [
  {
    id: PasswordRecommendationRuleIds.MIN_LENGTH,
    label: '8+ characters',
    test: (password: string) => password.length >= PasswordMinLength,
    validationMessage: 'Password must be at least 8 characters.'
  },
  {
    id: PasswordRecommendationRuleIds.NUMBER,
    label: 'One number',
    test: (password: string) => /[0-9]/u.test(password),
    validationMessage: 'Password must include at least one number.'
  },
  {
    id: PasswordRecommendationRuleIds.UPPERCASE,
    label: 'Uppercase letter',
    test: (password: string) => /[A-Z]/u.test(password),
    validationMessage: 'Password must include at least one uppercase letter.'
  },
  {
    id: PasswordRecommendationRuleIds.SPECIAL_CHARACTER,
    label: 'Special symbol',
    test: (password: string) => /[^A-Za-z0-9]/u.test(password),
    validationMessage: 'Password must include at least one special character.'
  }
] as const;
