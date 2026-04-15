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

export const appConfig = {
  appName: import.meta.env.VITE_APP_NAME ?? 'Estimate Room UI',
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
  environment: import.meta.env.MODE,
  oauthClientId: import.meta.env.VITE_OAUTH_CLIENT_ID ?? '',
  oauthRedirectUri:
    import.meta.env.VITE_OAUTH_REDIRECT_URI ?? resolveDefaultOauthRedirectUri(),
  oauthScopes: import.meta.env.VITE_OAUTH_SCOPES ?? 'openid user',
  wsBaseUrl: import.meta.env.VITE_WS_BASE_URL ?? resolveDefaultWsBaseUrl()
} as const;
