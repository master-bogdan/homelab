import { AUTH_STATUSES } from '../types';
import type { AuthState } from '../types';
import { authStateKey } from './auth.store';

type AuthStateRoot = {
  readonly [authStateKey]: AuthState;
};

export const selectAuthState = (state: AuthStateRoot) => state[authStateKey];
export const selectAuthStatus = (state: AuthStateRoot) => selectAuthState(state).status;
export const selectAuthUser = (state: AuthStateRoot) => selectAuthState(state).user;
export const selectIsAuthenticated = (state: AuthStateRoot) =>
  selectAuthStatus(state) === AUTH_STATUSES.AUTHENTICATED && selectAuthUser(state) !== null;
