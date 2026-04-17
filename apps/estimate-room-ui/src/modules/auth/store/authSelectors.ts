import { AuthStates } from '../constants';
import type { AuthState } from '../types';
import { AUTH_STATE_KEY } from './authStore';

type AuthStateRoot = {
  readonly [AUTH_STATE_KEY]: AuthState;
};

export const selectAuthState = (state: AuthStateRoot) => state[AUTH_STATE_KEY];
export const selectOAuthCallbackState = (state: AuthStateRoot) =>
  selectAuthState(state).oauthCallback;
export const selectAuthStatus = (state: AuthStateRoot) => selectAuthState(state).status;
export const selectAuthUser = (state: AuthStateRoot) => selectAuthState(state).user;
export const selectIsAuthenticated = (state: AuthStateRoot) =>
  selectAuthStatus(state) === AuthStates.AUTHENTICATED && selectAuthUser(state) !== null;
