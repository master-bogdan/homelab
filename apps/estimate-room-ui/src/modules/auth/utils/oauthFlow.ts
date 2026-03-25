import { appConfig } from '@/shared/config/env';

import type { PendingAuthorizationRequest } from '../types';

import { createApiUrl } from './apiUrl';

const pendingAuthorizationStorageKey = 'estimate-room.auth.pending-authorization';

const encodeBase64Url = (value: ArrayBuffer | Uint8Array) => {
  const bytes = value instanceof Uint8Array ? value : new Uint8Array(value);
  let binary = '';

  bytes.forEach((byte) => {
    binary += String.fromCharCode(byte);
  });

  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/u, '');
};

const createRandomToken = () => {
  const bytes = new Uint8Array(32);

  crypto.getRandomValues(bytes);

  return encodeBase64Url(bytes);
};

const createCodeChallenge = async (codeVerifier: string) => {
  const digest = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(codeVerifier));

  return encodeBase64Url(digest);
};

const getOauthConfig = () => {
  const clientId = appConfig.oauthClientId.trim();
  const redirectUri = appConfig.oauthRedirectUri.trim();
  const scopes = appConfig.oauthScopes.trim() || 'openid user';

  if (!clientId) {
    throw new Error('OAuth client is not configured for EstimateRoom UI.');
  }

  if (!redirectUri) {
    throw new Error('OAuth redirect URI is not configured for EstimateRoom UI.');
  }

  return { clientId, redirectUri, scopes };
};

export const readPendingAuthorizationRequest = (): PendingAuthorizationRequest | null => {
  if (typeof window === 'undefined') {
    return null;
  }

  const rawValue = window.sessionStorage.getItem(pendingAuthorizationStorageKey);

  if (!rawValue) {
    return null;
  }

  try {
    return JSON.parse(rawValue) as PendingAuthorizationRequest;
  } catch {
    window.sessionStorage.removeItem(pendingAuthorizationStorageKey);

    return null;
  }
};

export const clearPendingAuthorizationRequest = () => {
  if (typeof window === 'undefined') {
    return;
  }

  window.sessionStorage.removeItem(pendingAuthorizationStorageKey);
};

export const createPendingAuthorizationRequest = async (redirectTo: string) => {
  const { clientId, redirectUri, scopes } = getOauthConfig();
  const codeVerifier = createRandomToken();
  const codeChallenge = await createCodeChallenge(codeVerifier);
  const nonce = createRandomToken();
  const state = createRandomToken();
  const continueUrl = createApiUrl('oauth2/authorize', {
    client_id: clientId,
    code_challenge: codeChallenge,
    code_challenge_method: 'S256',
    nonce,
    redirect_uri: redirectUri,
    response_type: 'code',
    scopes,
    state
  }).toString();

  const request: PendingAuthorizationRequest = {
    clientId,
    codeVerifier,
    continueUrl,
    redirectTo,
    redirectUri,
    state
  };

  if (typeof window !== 'undefined') {
    window.sessionStorage.setItem(
      pendingAuthorizationStorageKey,
      JSON.stringify(request)
    );
  }

  return request;
};

export const ensurePendingAuthorizationRequest = async (
  redirectTo: string,
  continueUrl?: string | null
) => {
  if (continueUrl) {
    const request = readPendingAuthorizationRequest();

    if (!request || request.continueUrl !== continueUrl) {
      throw new Error('Your sign-in session expired. Please start again.');
    }

    return request;
  }

  return createPendingAuthorizationRequest(redirectTo);
};
