import { appConfig } from '@/shared/config/env';
import type { ApiError } from '@/shared/types';

type QueryValue = boolean | number | string | null | undefined;

export interface RequestOptions extends Omit<RequestInit, 'body'> {
  readonly body?: BodyInit | object;
  readonly query?: Record<string, QueryValue>;
}

const isBodyInit = (value: RequestOptions['body']): value is BodyInit =>
  value instanceof Blob ||
  value instanceof FormData ||
  value instanceof URLSearchParams ||
  typeof value === 'string' ||
  value instanceof ArrayBuffer;

const resolveBaseUrl = (baseUrl: string) => {
  if (/^https?:\/\//.test(baseUrl)) {
    return baseUrl;
  }

  const origin = typeof window === 'undefined' ? 'http://localhost' : window.location.origin;
  const normalizedBaseUrl = baseUrl.startsWith('/') ? baseUrl : `/${baseUrl}`;

  return `${origin}${normalizedBaseUrl}`;
};

const createRequestUrl = (
  baseUrl: string,
  path: string,
  query?: Record<string, QueryValue>
) => {
  const url = new URL(path.replace(/^\//, ''), `${resolveBaseUrl(baseUrl).replace(/\/$/, '')}/`);

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

const parseError = async (response: Response): Promise<ApiError> => {
  const fallbackError: ApiError = {
    message: `Request failed with status ${response.status}.`,
    status: response.status
  };

  const contentType = response.headers.get('content-type') ?? '';

  if (!contentType.includes('application/json')) {
    return fallbackError;
  }

  const payload = (await response.json()) as Partial<ApiError>;

  return {
    code: payload.code,
    details: payload.details,
    message: payload.message ?? fallbackError.message,
    status: payload.status ?? fallbackError.status
  };
};

export class HttpClient {
  public constructor(private readonly baseUrl: string) {}

  public async delete<TResponse>(path: string, options?: RequestOptions) {
    return this.request<TResponse>(path, {
      ...options,
      method: 'DELETE'
    });
  }

  public async get<TResponse>(path: string, options?: RequestOptions) {
    return this.request<TResponse>(path, {
      ...options,
      method: 'GET'
    });
  }

  public async post<TResponse>(path: string, body?: RequestOptions['body'], options?: RequestOptions) {
    return this.request<TResponse>(path, {
      ...options,
      body,
      method: 'POST'
    });
  }

  public async put<TResponse>(path: string, body?: RequestOptions['body'], options?: RequestOptions) {
    return this.request<TResponse>(path, {
      ...options,
      body,
      method: 'PUT'
    });
  }

  public async request<TResponse>(path: string, options: RequestOptions = {}): Promise<TResponse> {
    const { body, headers, query, ...restOptions } = options;
    const requestUrl = createRequestUrl(this.baseUrl, path, query);
    const requestHeaders = new Headers(headers);

    let requestBody: BodyInit | undefined;

    if (body && isBodyInit(body)) {
      requestBody = body;
    } else if (body) {
      requestBody = JSON.stringify(body);
      requestHeaders.set('Content-Type', 'application/json');
    }

    const response = await fetch(requestUrl, {
      credentials: 'include',
      ...restOptions,
      body: requestBody,
      headers: requestHeaders
    });

    if (!response.ok) {
      throw await parseError(response);
    }

    if (response.status === 204) {
      return undefined as TResponse;
    }

    const contentType = response.headers.get('content-type') ?? '';

    if (contentType.includes('application/json')) {
      return (await response.json()) as TResponse;
    }

    return (await response.text()) as TResponse;
  }
}

export const apiClient = new HttpClient(appConfig.apiBaseUrl);
