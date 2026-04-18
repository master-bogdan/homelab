export const ResetPasswordErrorMessages = {
  RESET_FAILED: 'Unable to reset your password right now.',
  TOKEN_VALIDATION_FAILED: 'Unable to validate this reset link.'
} as const;

export const ResetPasswordSearchParams = {
  TOKEN: 'token'
} as const;

export const ResetPasswordValidationMessageFragments = {
  EXPIRED: 'expired reset token',
  INVALID: 'invalid reset token',
  USED: 'used reset token'
} as const;

export const ResetPasswordValidationReasons = {
  EXPIRED: 'expired',
  INVALID: 'invalid',
  USED: 'used'
} as const;

export const ResetPasswordInvalidLinkCopy = {
  [ResetPasswordValidationReasons.EXPIRED]: {
    description: 'This password reset link has expired. Request a new link to continue.',
    title: 'Expired Link'
  },
  [ResetPasswordValidationReasons.INVALID]: {
    description:
      'This password reset link is invalid or has expired. Please request a new link to reset your password.',
    title: 'Invalid Link'
  },
  [ResetPasswordValidationReasons.USED]: {
    description:
      'This password reset link has already been used. Request a new link to set another password.',
    title: 'Link Already Used'
  }
} as const;
