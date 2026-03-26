import { appConfig } from '@/shared/config/env';

type QueryValue = boolean | number | string | null | undefined;

const resolveBaseUrl = (baseUrl: string) => {
  if (/^https?:\/\//.test(baseUrl)) {
    return baseUrl;
  }

  const origin = typeof window === 'undefined' ? 'http://localhost' : window.location.origin;
  const normalizedBaseUrl = baseUrl.startsWith('/') ? baseUrl : `/${baseUrl}`;

  return `${origin}${normalizedBaseUrl}`;
};

export const createApiUrl = (path: string, query?: Record<string, QueryValue>) => {
  const url = new URL(
    path.replace(/^\//, ''),
    `${resolveBaseUrl(appConfig.apiBaseUrl).replace(/\/$/, '')}/`
  );

  if (query) {
    Object.entries(query).forEach(([key, value]) => {
      if (value === undefined || value === null) {
        return;
      }

      url.searchParams.set(key, String(value));
    });
  }

  return url;
};
