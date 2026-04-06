describe('oauthFlow', () => {
  afterEach(() => {
    sessionStorage.clear();
    vi.resetModules();
    vi.unstubAllEnvs();
  });

  it('creates and stores a PKCE authorization request for the backend continue flow', async () => {
    vi.stubEnv('VITE_OAUTH_CLIENT_ID', 'estimate-room-ui');
    vi.stubEnv('VITE_OAUTH_REDIRECT_URI', 'http://localhost:5173/auth/callback');
    vi.stubEnv('VITE_OAUTH_SCOPES', 'openid user');

    const {
      createPendingAuthorizationRequest,
      readPendingAuthorizationRequest
    } = await import('../oauthFlow');

    const request = await createPendingAuthorizationRequest('/dashboard?tab=activity');
    const continueUrl = new URL(request.continueUrl, 'http://localhost');

    expect(request.continueUrl.startsWith('/api/v1/oauth2/authorize?')).toBe(true);
    expect(continueUrl.pathname).toBe('/api/v1/oauth2/authorize');
    expect(continueUrl.searchParams.get('client_id')).toBe('estimate-room-ui');
    expect(continueUrl.searchParams.get('redirect_uri')).toBe(
      'http://localhost:5173/auth/callback'
    );
    expect(continueUrl.searchParams.get('response_type')).toBe('code');
    expect(continueUrl.searchParams.get('scopes')).toBe('openid user');
    expect(continueUrl.searchParams.get('code_challenge')).toBeTruthy();
    expect(continueUrl.searchParams.get('code_challenge_method')).toBe('S256');
    expect(continueUrl.searchParams.get('state')).toBe(request.state);
    expect(request.redirectTo).toBe('/dashboard?tab=activity');
    expect(readPendingAuthorizationRequest()).toEqual(request);
  });

  it('fails when a backend continue URL exists without a matching stored transaction', async () => {
    const { ensurePendingAuthorizationRequest } = await import('../oauthFlow');

    await expect(
      ensurePendingAuthorizationRequest(
        '/dashboard',
        '/api/v1/oauth2/authorize?client_id=client&redirect_uri=http%3A%2F%2Flocalhost%3A5173%2Fauth%2Fcallback&response_type=code&scopes=openid%20user&state=test&code_challenge=test&code_challenge_method=S256&nonce=test'
      )
    ).rejects.toThrow('Your sign-in session expired. Please start again.');
  });
});
