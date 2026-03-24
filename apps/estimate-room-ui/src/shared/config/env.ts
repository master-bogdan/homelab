const resolveDefaultWsBaseUrl = () => {
  if (typeof window === 'undefined') {
    return 'ws://localhost:8080/ws';
  }

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';

  return `${protocol}//${window.location.host}/ws`;
};

export const appConfig = {
  appName: import.meta.env.VITE_APP_NAME ?? 'Estimate Room UI',
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL ?? '/api',
  environment: import.meta.env.MODE,
  wsBaseUrl: import.meta.env.VITE_WS_BASE_URL ?? resolveDefaultWsBaseUrl()
} as const;
