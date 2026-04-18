import { render, screen, waitFor } from '@testing-library/react';
import { StrictMode } from 'react';
import { Provider } from 'react-redux';
import { MemoryRouter, Route, Routes } from 'react-router-dom';

import { createTestStore } from '@/test/test-utils';

import { AuthStates } from '../../constants';
import { useOAuthCallbackPage } from '../useOAuthCallbackPage';

const createJsonResponse = (payload: unknown, status = 200) =>
  new Response(JSON.stringify(payload), {
    headers: {
      'content-type': 'application/json'
    },
    status
  });

const OAuthCallbackRoute = () => {
  const { errorMessage, isLoading } = useOAuthCallbackPage();

  return <div>{isLoading ? 'Loading callback' : errorMessage}</div>;
};

describe('useOAuthCallbackPage', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    sessionStorage.clear();
    fetchMock = vi.fn();
    vi.stubGlobal('fetch', fetchMock);
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
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
          displayName: 'Test User',
          email: 'test@example.com',
          id: 'user-1',
          occupation: null,
          organization: null
        }
      }));

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
      expect(fetchMock).toHaveBeenCalledTimes(2);
    });
    await screen.findByText('Dashboard');
    await waitFor(() => {
      expect(store.getState().auth.status).toBe(AuthStates.AUTHENTICATED);
    });

    expect(sessionStorage.getItem('estimate-room.auth.pending-authorization')).toBeNull();
  });
});
