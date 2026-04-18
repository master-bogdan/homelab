const resolveDefaultWsBaseUrl = () => {
  if (typeof window === 'undefined') {
    return 'ws://localhost:8080/ws';
  }

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';

  return `${protocol}//${window.location.host}/ws`;
};

const resolveDefaultOauthRedirectUri = () => {
  if (typeof window === 'undefined') {
    return 'http://localhost:5173/auth/callback';
  }

  return `${window.location.origin}/auth/callback`;
};

export const AppConfig = {
  APP_NAME: import.meta.env.VITE_APP_NAME ?? 'Estimate Room UI',
  API_BASE_URL: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
  ENVIRONMENT: import.meta.env.MODE,
  OAUTH_CLIENT_ID: import.meta.env.VITE_OAUTH_CLIENT_ID ?? '',
  OAUTH_REDIRECT_URI:
    import.meta.env.VITE_OAUTH_REDIRECT_URI ?? resolveDefaultOauthRedirectUri(),
  OAUTH_SCOPES: import.meta.env.VITE_OAUTH_SCOPES ?? 'openid user',
  WS_BASE_URL: import.meta.env.VITE_WS_BASE_URL ?? resolveDefaultWsBaseUrl()
} as const;
