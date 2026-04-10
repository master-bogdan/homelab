import type { ApiError } from '@/shared/types';

import type { ResetPasswordValidationReason } from '../types';

type RtkQueryError = {
  readonly data?: unknown;
  readonly error?: string;
  readonly status: number | string;
};

const isApiError = (error: unknown): error is ApiError =>
  typeof error === 'object' &&
  error !== null &&
  'status' in error &&
  typeof (error as { status?: unknown }).status === 'number';

const isRtkQueryError = (error: unknown): error is RtkQueryError =>
  typeof error === 'object' &&
  error !== null &&
  'status' in error &&
  (typeof (error as { status?: unknown }).status === 'number' ||
    typeof (error as { status?: unknown }).status === 'string');

const normalizeErrorText = (value: string) => value.trim().toLowerCase();

const extractRtkQueryMessage = (error: RtkQueryError) => {
  if (
    error.data &&
    typeof error.data === 'object' &&
    'detail' in error.data &&
    typeof (error.data as { detail?: unknown }).detail === 'string'
  ) {
    return (error.data as { detail: string }).detail;
  }

  if (
    error.data &&
    typeof error.data === 'object' &&
    'message' in error.data &&
    typeof (error.data as { message?: unknown }).message === 'string'
  ) {
    return (error.data as { message: string }).message;
  }

  if (
    error.data &&
    typeof error.data === 'object' &&
    'title' in error.data &&
    typeof (error.data as { title?: unknown }).title === 'string'
  ) {
    return (error.data as { title: string }).title;
  }

  return error.error ?? '';
};

const extractErrorText = (error: unknown) => {
  if (isApiError(error)) {
    return normalizeErrorText(error.detail ?? error.message ?? error.title ?? '');
  }

  if (isRtkQueryError(error)) {
    return normalizeErrorText(extractRtkQueryMessage(error));
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

  if (isRtkQueryError(error)) {
    return extractRtkQueryMessage(error) || fallbackMessage;
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
