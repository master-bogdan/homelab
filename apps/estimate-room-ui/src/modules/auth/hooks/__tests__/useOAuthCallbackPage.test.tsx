import { render, screen, waitFor } from '@testing-library/react';
import { StrictMode } from 'react';
import { Provider } from 'react-redux';
import { MemoryRouter, Route, Routes } from 'react-router-dom';

import { createTestStore } from '@/test/test-utils';

import { authService } from '../../services';
import { useOAuthCallbackPage } from '../useOAuthCallbackPage';

vi.mock('@/modules/auth/services', () => ({
  authService: {
    exchangeAuthorizationCode: vi.fn(),
    fetchSession: vi.fn()
  }
}));

const mockedAuthService = vi.mocked(authService);

const OAuthCallbackRoute = () => {
  const { errorMessage, isLoading } = useOAuthCallbackPage();

  return <div>{isLoading ? 'Loading callback' : errorMessage}</div>;
};

describe('useOAuthCallbackPage', () => {
  beforeEach(() => {
    sessionStorage.clear();
    vi.clearAllMocks();
  });

  it('exchanges the authorization code only once in StrictMode', async () => {
    sessionStorage.setItem(
      'estimate-room.auth.pending-authorization',
      JSON.stringify({
        clientId: 'estimate-room-ui',
        codeVerifier: 'verifier-123',
        continueUrl: '/api/v1/oauth2/authorize?...',
        redirectTo: '/dashboard',
        redirectUri: 'http://localhost:5173/auth/callback',
        state: 'state-123'
      })
    );

    mockedAuthService.exchangeAuthorizationCode.mockResolvedValue({
      accessToken: 'access-token',
      expiresIn: 900,
      idToken: 'id-token',
      refreshToken: 'refresh-token',
      tokenType: 'Bearer'
    });
    mockedAuthService.fetchSession.mockResolvedValue({
      avatarUrl: null,
      displayName: 'Test User',
      email: 'test@example.com',
      id: 'user-1',
      occupation: null,
      organization: null
    });

    const store = createTestStore();
    render(
      <StrictMode>
        <Provider store={store}>
          <MemoryRouter initialEntries={['/auth/callback?code=code-123&state=state-123']}>
            <Routes>
              <Route element={<OAuthCallbackRoute />} path="/auth/callback" />
              <Route element={<div>Dashboard</div>} path="/dashboard" />
            </Routes>
          </MemoryRouter>
        </Provider>
      </StrictMode>
    );

    await waitFor(() => {
      expect(mockedAuthService.exchangeAuthorizationCode).toHaveBeenCalledTimes(1);
    });
    await waitFor(() => {
      expect(mockedAuthService.fetchSession).toHaveBeenCalledTimes(1);
    });
    await screen.findByText('Dashboard');
    await waitFor(() => {
      expect(store.getState().auth.status).toBe('authenticated');
    });

    expect(sessionStorage.getItem('estimate-room.auth.pending-authorization')).toBeNull();
  });
});
