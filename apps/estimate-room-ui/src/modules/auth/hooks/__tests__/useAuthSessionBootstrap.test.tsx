import { act, renderHook, waitFor } from '@testing-library/react';
import type { PropsWithChildren } from 'react';
import { StrictMode } from 'react';
import { Provider } from 'react-redux';

import { accessTokenStorage } from '@/shared/api';
import { createTestStore } from '@/test/test-utils';

import { AuthStates } from '../../types';
import { useAuthSessionBootstrap } from '../useAuthSessionBootstrap';

const createJsonResponse = (payload: unknown, status = 200) =>
  new Response(JSON.stringify(payload), {
    headers: {
      'content-type': 'application/json'
    },
    status
  });

describe('useAuthSessionBootstrap', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    window.history.replaceState({}, '', '/login');
    window.localStorage.clear();
    fetchMock = vi.fn();
    vi.stubGlobal('fetch', fetchMock);
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  it('fetches the session only once in StrictMode on non-callback routes', async () => {
    accessTokenStorage.set('access-token');
    fetchMock.mockResolvedValue(createJsonResponse({
      authenticated: false,
      user: null
    }));

    const store = createTestStore();
    const wrapper = ({ children }: PropsWithChildren) => (
      <StrictMode>
        <Provider store={store}>{children}</Provider>
      </StrictMode>
    );

    renderHook(() => useAuthSessionBootstrap(), { wrapper });

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(1);
    });
    await waitFor(() => {
      expect(store.getState().auth.status).toBe(AuthStates.UNAUTHENTICATED);
    });
  });

  it('skips the bootstrap request on the OAuth callback route', async () => {
    window.history.replaceState({}, '', '/auth/callback?code=code-123&state=state-123');

    const store = createTestStore();
    const wrapper = ({ children }: PropsWithChildren) => (
      <StrictMode>
        <Provider store={store}>{children}</Provider>
      </StrictMode>
    );

    renderHook(() => useAuthSessionBootstrap(), { wrapper });

    await act(async () => {
      await Promise.resolve();
    });

    expect(fetchMock).not.toHaveBeenCalled();
    expect(store.getState().auth.status).toBe(AuthStates.UNKNOWN);
  });

  it('refreshes the access token before fetching the session when local storage is empty', async () => {
    fetchMock
      .mockResolvedValueOnce(createJsonResponse({
        access_token: 'access-token',
        expires_in: 900,
        id_token: 'id-token',
        token_type: 'Bearer'
      }))
      .mockResolvedValueOnce(createJsonResponse({
        authenticated: true,
        user: {
          avatarUrl: null,
          displayName: 'Alex Architect',
          email: 'alex@example.com',
          id: 'user-1',
          occupation: null,
          organization: null
        }
      }));

    const store = createTestStore();
    const wrapper = ({ children }: PropsWithChildren) => (
      <StrictMode>
        <Provider store={store}>{children}</Provider>
      </StrictMode>
    );

    renderHook(() => useAuthSessionBootstrap(), { wrapper });

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(2);
    });
    await waitFor(() => {
      expect(store.getState().auth.status).toBe(AuthStates.AUTHENTICATED);
    });
  });
});
