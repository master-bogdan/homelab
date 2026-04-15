import type { ApiError } from '@/shared/types';
import { appRoutes } from '@/shared/constants/routes';

import type {
  DashboardDeckPreset,
  DashboardDeckPresetKey,
  DashboardSessionStatus
} from '../types';

type RtkQueryError = {
  readonly data?: unknown;
  readonly error?: string;
  readonly status: number | string;
};

const relativeTimeFormatter = new Intl.RelativeTimeFormat('en-US', {
  numeric: 'auto'
});

const dateTimeFormatter = new Intl.DateTimeFormat('en-US', {
  day: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
  month: 'short'
});

const dashboardDeckPresetMap: Record<DashboardDeckPresetKey, DashboardDeckPreset> = {
  fibonacci: {
    deck: {
      kind: 'FIBONACCI',
      name: 'Fibonacci',
      values: ['0', '1', '2', '3', '5', '8', '13', '21', '?']
    },
    description: 'Default planning poker scale used across EstimateRoom.',
    key: 'fibonacci',
    label: 'Fibonacci'
  },
  powerOfTwo: {
    deck: {
      kind: 'POWER_OF_TWO',
      name: 'Power of 2',
      values: ['1', '2', '4', '8', '16', '32', '?']
    },
    description: 'A wider step curve for larger infrastructure sizing conversations.',
    key: 'powerOfTwo',
    label: 'Power of 2'
  },
  simple: {
    deck: {
      kind: 'SIMPLE',
      name: 'Simple 1-5',
      values: ['1', '2', '3', '4', '5']
    },
    description: 'Compact numeric range for lightweight estimation passes.',
    key: 'simple',
    label: 'Simple 1-5'
  },
  tShirt: {
    deck: {
      kind: 'T_SHIRT',
      name: 'T-Shirt',
      values: ['XS', 'S', 'M', 'L', 'XL', '?']
    },
    description: 'Relative sizing for discovery work and early architectural shaping.',
    key: 'tShirt',
    label: 'T-Shirt'
  }
};

export const dashboardDeckPresets = Object.values(dashboardDeckPresetMap);

export const getDashboardDeckPreset = (key: DashboardDeckPresetKey) =>
  dashboardDeckPresetMap[key];

export const parseInviteEmails = (value: string) =>
  Array.from(
    new Set(
      value
        .split(/[\n,;]+/)
        .map((entry) => entry.trim())
        .filter(Boolean)
    )
  );

export const extractInviteToken = (value: string) => {
  const trimmedValue = value.trim();

  if (!trimmedValue) {
    return '';
  }

  try {
    const parsedUrl = new URL(trimmedValue);
    const queryToken =
      parsedUrl.searchParams.get('token') ??
      parsedUrl.searchParams.get('inviteToken') ??
      parsedUrl.searchParams.get('code');

    if (queryToken) {
      return queryToken.trim();
    }

    const pathSegments = parsedUrl.pathname.split('/').filter(Boolean);

    return pathSegments.at(-1)?.trim() ?? trimmedValue;
  } catch {
    return trimmedValue.replace(/\/+$/, '');
  }
};

export const formatDashboardDateTime = (value: string) => dateTimeFormatter.format(new Date(value));

export const buildDashboardInviteLink = (roomCode: string) => {
  const origin = typeof window === 'undefined' ? 'http://localhost:5173' : window.location.origin;

  return `${origin}${appRoutes.joinRoomPath(roomCode)}`;
};

export const formatRelativeTime = (value: string, now = new Date()) => {
  const deltaSeconds = Math.round((new Date(value).getTime() - now.getTime()) / 1000);
  const absoluteSeconds = Math.abs(deltaSeconds);

  if (absoluteSeconds < 60) {
    return relativeTimeFormatter.format(Math.round(deltaSeconds), 'second');
  }

  const deltaMinutes = Math.round(deltaSeconds / 60);
  if (Math.abs(deltaMinutes) < 60) {
    return relativeTimeFormatter.format(deltaMinutes, 'minute');
  }

  const deltaHours = Math.round(deltaMinutes / 60);
  if (Math.abs(deltaHours) < 24) {
    return relativeTimeFormatter.format(deltaHours, 'hour');
  }

  const deltaDays = Math.round(deltaHours / 24);
  if (Math.abs(deltaDays) < 30) {
    return relativeTimeFormatter.format(deltaDays, 'day');
  }

  const deltaMonths = Math.round(deltaDays / 30);
  if (Math.abs(deltaMonths) < 12) {
    return relativeTimeFormatter.format(deltaMonths, 'month');
  }

  return relativeTimeFormatter.format(Math.round(deltaMonths / 12), 'year');
};

export const formatStatusLabel = (value: string) =>
  value
    .toLowerCase()
    .split('_')
    .map((segment) => segment.charAt(0).toUpperCase() + segment.slice(1))
    .join(' ');

const isRtkQueryError = (error: unknown): error is RtkQueryError =>
  typeof error === 'object' &&
  error !== null &&
  'status' in error &&
  (typeof (error as { status?: unknown }).status === 'number' ||
    typeof (error as { status?: unknown }).status === 'string');

const getRtkQueryErrorMessage = (error: RtkQueryError) => {
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

export const getDashboardErrorMessage = (error: unknown, fallback: string) => {
  if (typeof error === 'string' && error.trim()) {
    return error;
  }

  if (!error || typeof error !== 'object') {
    return fallback;
  }

  const apiError = error as Partial<ApiError>;

  if (apiError.detail || apiError.message) {
    return apiError.detail ?? apiError.message ?? fallback;
  }

  if (isRtkQueryError(error)) {
    return getRtkQueryErrorMessage(error) || fallback;
  }

  return fallback;
};

export const getInitials = (value: string) =>
  value
    .split(' ')
    .map((segment) => segment.trim().charAt(0))
    .filter(Boolean)
    .slice(0, 2)
    .join('')
    .toUpperCase();

export const getSessionDestinationLabel = (status: DashboardSessionStatus) =>
  status === 'ACTIVE' ? 'Open room' : 'Review summary';

export const getArchitectLevelLabel = (level: number) => `Level ${level} Architect`;

export const getXpHint = (xp: number, nextLevelXp: number, level: number) => {
  if (nextLevelXp <= xp) {
    return `Level ${level} is fully charged.`;
  }

  const remainingXp = nextLevelXp - xp;

  return `Earn ${remainingXp} more XP to reach level ${level + 1}.`;
};
