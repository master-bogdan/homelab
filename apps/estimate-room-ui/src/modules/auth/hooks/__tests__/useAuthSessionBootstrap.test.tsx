import { act, renderHook, waitFor } from '@testing-library/react';
import type { PropsWithChildren } from 'react';
import { StrictMode } from 'react';
import { Provider } from 'react-redux';

import { createTestStore } from '@/test/test-utils';

import { authService } from '../../services';
import { AUTH_STATUSES } from '../../types';
import { useAuthSessionBootstrap } from '../useAuthSessionBootstrap';

vi.mock('@/modules/auth/services', () => ({
  authService: {
    fetchSession: vi.fn(),
    hasStoredAccessToken: vi.fn(),
    logout: vi.fn(),
    refreshAccessToken: vi.fn()
  }
}));

const mockedAuthService = vi.mocked(authService);

describe('useAuthSessionBootstrap', () => {
  beforeEach(() => {
    window.history.replaceState({}, '', '/login');
    vi.clearAllMocks();
    mockedAuthService.hasStoredAccessToken.mockReturnValue(true);
    mockedAuthService.logout.mockResolvedValue({ loggedOut: true });
    mockedAuthService.refreshAccessToken.mockResolvedValue({
      accessToken: 'access-token',
      expiresIn: 900,
      idToken: 'id-token',
      refreshToken: 'refresh-token',
      tokenType: 'Bearer'
    });
  });

  it('fetches the session only once in StrictMode on non-callback routes', async () => {
    mockedAuthService.fetchSession.mockResolvedValue(null);

    const store = createTestStore();
    const wrapper = ({ children }: PropsWithChildren) => (
      <StrictMode>
        <Provider store={store}>{children}</Provider>
      </StrictMode>
    );

    renderHook(() => useAuthSessionBootstrap(), { wrapper });

    await waitFor(() => {
      expect(mockedAuthService.fetchSession).toHaveBeenCalledTimes(1);
    });
    await waitFor(() => {
      expect(store.getState().auth.status).toBe(AUTH_STATUSES.UNAUTHENTICATED);
    });
  });

  it('skips the bootstrap request on the OAuth callback route', async () => {
    window.history.replaceState({}, '', '/auth/callback?code=code-123&state=state-123');
    mockedAuthService.fetchSession.mockResolvedValue(null);

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

    expect(mockedAuthService.fetchSession).not.toHaveBeenCalled();
    expect(store.getState().auth.status).toBe(AUTH_STATUSES.UNKNOWN);
  });

  it('refreshes the access token before fetching the session when local storage is empty', async () => {
    mockedAuthService.hasStoredAccessToken.mockReturnValue(false);
    mockedAuthService.fetchSession.mockResolvedValue({
      avatarUrl: null,
      displayName: 'Alex Architect',
      email: 'alex@example.com',
      id: 'user-1',
      occupation: null,
      organization: null
    });

    const store = createTestStore();
    const wrapper = ({ children }: PropsWithChildren) => (
      <StrictMode>
        <Provider store={store}>{children}</Provider>
      </StrictMode>
    );

    renderHook(() => useAuthSessionBootstrap(), { wrapper });

    await waitFor(() => {
      expect(mockedAuthService.refreshAccessToken).toHaveBeenCalledTimes(1);
    });
    await waitFor(() => {
      expect(mockedAuthService.fetchSession).toHaveBeenCalledTimes(1);
    });
    await waitFor(() => {
      expect(store.getState().auth.status).toBe(AUTH_STATUSES.AUTHENTICATED);
    });
  });
});
