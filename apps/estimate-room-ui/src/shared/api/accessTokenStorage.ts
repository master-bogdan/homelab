const ACCESS_TOKEN_STORAGE_KEY = 'estimate-room.auth.access-token';

export const accessTokenStorage = {
  clear() {
    if (typeof window === 'undefined') {
      return;
    }

    window.localStorage.removeItem(ACCESS_TOKEN_STORAGE_KEY);
  },
  get() {
    if (typeof window === 'undefined') {
      return null;
    }

    return window.localStorage.getItem(ACCESS_TOKEN_STORAGE_KEY);
  },
  set(accessToken: string) {
    if (typeof window === 'undefined') {
      return;
    }

    window.localStorage.setItem(ACCESS_TOKEN_STORAGE_KEY, accessToken);
  }
};
