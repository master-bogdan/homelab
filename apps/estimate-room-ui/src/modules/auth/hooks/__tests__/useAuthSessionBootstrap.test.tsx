import { act, renderHook, waitFor } from '@testing-library/react';
import type { PropsWithChildren } from 'react';
import { StrictMode } from 'react';
import { Provider } from 'react-redux';

import { createTestStore } from '@/test/test-utils';

import { authService } from '../../services';
import { useAuthSessionBootstrap } from '../useAuthSessionBootstrap';

vi.mock('@/modules/auth/services', () => ({
  authService: {
    fetchSession: vi.fn()
  }
}));

const mockedAuthService = vi.mocked(authService);

describe('useAuthSessionBootstrap', () => {
  beforeEach(() => {
    window.history.replaceState({}, '', '/login');
    vi.clearAllMocks();
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
      expect(store.getState().auth.status).toBe('unauthenticated');
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
    expect(store.getState().auth.status).toBe('unknown');
  });
});
