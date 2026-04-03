import type { OAuthTokenResponse } from '../types';

const createCookieAttributes = () => {
  const secureAttribute =
    typeof window !== 'undefined' && window.location.protocol === 'https:' ? '; Secure' : '';

  return `Path=/; SameSite=Lax${secureAttribute}`;
};

const writeCookie = (name: string, value: string, extraAttributes = '') => {
  if (typeof document === 'undefined') {
    return;
  }

  document.cookie = `${name}=${encodeURIComponent(value)}; ${createCookieAttributes()}${extraAttributes}`;
};

export const persistOauthTokenCookies = (tokens: OAuthTokenResponse) => {
  writeCookie('access_token', tokens.accessToken, `; Max-Age=${tokens.expiresIn}`);

  if (tokens.refreshToken) {
    writeCookie('refresh_token', tokens.refreshToken);
  }
};

export const clearOauthTokenCookies = () => {
  writeCookie('access_token', '', '; Max-Age=0');
  writeCookie('refresh_token', '', '; Max-Age=0');
};
