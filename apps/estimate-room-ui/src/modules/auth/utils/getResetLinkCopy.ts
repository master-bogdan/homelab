import type { ResetPasswordValidationReason } from '../types';

export const getResetLinkCopy = (reason: ResetPasswordValidationReason | undefined) => {
  switch (reason) {
    case 'expired':
      return {
        description:
          'This password reset link has expired. Request a new link to continue.',
        title: 'Expired Link'
      };
    case 'used':
      return {
        description:
          'This password reset link has already been used. Request a new link to set another password.',
        title: 'Link Already Used'
      };
    case 'invalid':
    default:
      return {
        description:
          'This password reset link is invalid or has expired. Please request a new link to reset your password.',
        title: 'Invalid Link'
      };
  }
};
