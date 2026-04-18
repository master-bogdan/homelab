import {
  fetchBaseQuery,
  type BaseQueryFn,
  type FetchArgs,
  type FetchBaseQueryError
} from '@reduxjs/toolkit/query/react';

import { AppConfig } from '@/config';

import { accessTokenStorage } from './accessTokenStorage';
import { apiSessionExpired } from './sessionLifecycle';

interface RefreshTokenApiResponse {
  readonly access_token?: string;
}

let refreshPromise: Promise<string | null> | null = null;

const rawBaseQuery = fetchBaseQuery({
  baseUrl: AppConfig.API_BASE_URL,
  credentials: 'include',
  prepareHeaders: (headers) => {
    const accessToken = accessTokenStorage.get();

    if (accessToken) {
      headers.set('authorization', `Bearer ${accessToken}`);
    }

    headers.set('accept', 'application/json');

    return headers;
  }
});

const normalizePath = (args: string | FetchArgs) => {
  const url = typeof args === 'string' ? args : args.url;

  return url.replace(/^\//, '');
};

const shouldRefresh = (args: string | FetchArgs) => {
  const path = normalizePath(args);

  return path === 'auth/session' || (!path.startsWith('auth/') && path !== 'oauth2/token');
};

async function refreshAccessToken(
  api: Parameters<BaseQueryFn>[1],
  extraOptions: Parameters<BaseQueryFn>[2]
): Promise<string | null> {
  if (!refreshPromise) {
    refreshPromise = (async () => {
      const clientId = AppConfig.OAUTH_CLIENT_ID.trim();

      if (!clientId) {
        accessTokenStorage.clear();
        api.dispatch(apiSessionExpired());

        return null;
      }

      const refreshResult = await rawBaseQuery(
        {
          body: new URLSearchParams({
            client_id: clientId,
            grant_type: 'refresh_token'
          }).toString(),
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
          },
          method: 'POST',
          url: 'oauth2/token'
        },
        api,
        extraOptions
      );

      if (refreshResult.data) {
        const data = refreshResult.data as RefreshTokenApiResponse;

        if (data.access_token) {
          accessTokenStorage.set(data.access_token);

          return data.access_token;
        }
      }

      accessTokenStorage.clear();
      api.dispatch(apiSessionExpired());
      await rawBaseQuery(
        {
          method: 'POST',
          url: 'auth/logout'
        },
        api,
        extraOptions
      );

      return null;
    })().finally(() => {
      refreshPromise = null;
    });
  }

  return refreshPromise;
}

export const baseQueryWithAuth: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError
> = async (args, api, extraOptions) => {
  let result = await rawBaseQuery(args, api, extraOptions);

  if (result.error?.status === 401 && shouldRefresh(args)) {
    const accessToken = await refreshAccessToken(api, extraOptions);

    if (accessToken) {
      result = await rawBaseQuery(args, api, extraOptions);
    }
  }

  return result;
};
