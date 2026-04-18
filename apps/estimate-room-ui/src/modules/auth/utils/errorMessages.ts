import type { ApiError } from '@/shared/types';

type RtkQueryError = {
  readonly data?: unknown;
  readonly error?: string;
  readonly status: number | string;
};

type ErrorPayload = {
  readonly detail?: unknown;
  readonly message?: unknown;
  readonly title?: unknown;
};

const ErrorMessageKeys = ['detail', 'message', 'title'] as const;

const isRecord = (value: unknown): value is Record<PropertyKey, unknown> =>
  typeof value === 'object' && value !== null;

const hasStatus = (value: unknown): value is { readonly status: unknown } =>
  isRecord(value) && 'status' in value;

const isApiError = (error: unknown): error is ApiError =>
  hasStatus(error) && typeof error.status === 'number';

const isRtkQueryError = (error: unknown): error is RtkQueryError =>
  hasStatus(error) &&
  (typeof error.status === 'string' || 'data' in error || 'error' in error);

const normalizeErrorText = (value: string) => value.trim().toLowerCase();

const getPayloadMessage = (payload: ErrorPayload) => {
  const messageKey = ErrorMessageKeys.find((key) => typeof payload[key] === 'string');

  return messageKey ? String(payload[messageKey]) : '';
};

const getApiErrorMessage = (error: ApiError) => getPayloadMessage(error);

const getRtkQueryErrorMessage = (error: RtkQueryError) => {
  if (isRecord(error.data)) {
    return getPayloadMessage(error.data);
  }

  return error.error ?? '';
};

const getErrorMessage = (error: unknown) => {
  if (isRtkQueryError(error)) {
    return getRtkQueryErrorMessage(error);
  }

  if (isApiError(error)) {
    return getApiErrorMessage(error);
  }

  if (error instanceof Error) {
    return error.message;
  }

  return '';
};

const extractErrorText = (error: unknown) => {
  const message = getErrorMessage(error);

  return message ? normalizeErrorText(message) : '';
};

export const resolveApiErrorMessage = (error: unknown, fallbackMessage: string) => {
  const message = getErrorMessage(error);

  return message || fallbackMessage;
};

export const isInvalidCredentialsError = (error: unknown) =>
  extractErrorText(error).includes('invalid credentials');

export const isEmailAlreadyInUseError = (error: unknown) =>
  extractErrorText(error).includes('email already in use');
