import { AppConfig } from '@/config';

type QueryValue = boolean | number | string | null | undefined;

const resolveBaseUrl = (baseUrl: string) => {
  if (/^https?:\/\//.test(baseUrl)) {
    return baseUrl;
  }

  const origin = typeof window === 'undefined' ? 'http://localhost' : window.location.origin;
  const normalizedBaseUrl = baseUrl.startsWith('/') ? baseUrl : `/${baseUrl}`;

  return `${origin}${normalizedBaseUrl}`;
};

const resolveOrigin = (baseUrl: string) => {
  if (/^https?:\/\//.test(baseUrl)) {
    return new URL(baseUrl).origin;
  }

  return typeof window === 'undefined' ? 'http://localhost' : window.location.origin;
};

export const createApiUrl = (path: string, query?: Record<string, QueryValue>) => {
  const url = new URL(
    path.replace(/^\//, ''),
    `${resolveBaseUrl(AppConfig.API_BASE_URL).replace(/\/$/, '')}/`
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

export const createApiPath = (path: string, query?: Record<string, QueryValue>) => {
  const url = createApiUrl(path, query);

  return `${url.pathname}${url.search}`;
};

export const resolveApiHref = (pathOrUrl: string) => {
  if (/^https?:\/\//.test(pathOrUrl)) {
    return pathOrUrl;
  }

  if (pathOrUrl.startsWith('/')) {
    return new URL(pathOrUrl, `${resolveOrigin(AppConfig.API_BASE_URL)}/`).toString();
  }

  return createApiUrl(pathOrUrl).toString();
};

export const createGithubLoginUrl = (continueUrl: string) =>
  createApiUrl('auth/github/login', { continue: continueUrl }).toString();
