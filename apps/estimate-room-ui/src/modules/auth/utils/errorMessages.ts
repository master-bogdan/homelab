import type { ApiError } from '@/shared/types';

import type { ResetPasswordValidationReason } from '../types';

const isApiError = (error: unknown): error is ApiError =>
  typeof error === 'object' &&
  error !== null &&
  'status' in error &&
  typeof (error as { status?: unknown }).status === 'number';

const normalizeErrorText = (value: string) => value.trim().toLowerCase();

const extractErrorText = (error: unknown) => {
  if (isApiError(error)) {
    return normalizeErrorText(error.detail ?? error.message ?? error.title ?? '');
  }

  if (error instanceof Error) {
    return normalizeErrorText(error.message);
  }

  return '';
};

export const resolveApiErrorMessage = (error: unknown, fallbackMessage: string) => {
  if (isApiError(error)) {
    return error.detail ?? error.message ?? error.title ?? fallbackMessage;
  }

  if (error instanceof Error && error.message) {
    return error.message;
  }

  return fallbackMessage;
};

export const isInvalidCredentialsError = (error: unknown) =>
  extractErrorText(error).includes('invalid credentials');

export const isEmailAlreadyInUseError = (error: unknown) =>
  extractErrorText(error).includes('email already in use');

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
