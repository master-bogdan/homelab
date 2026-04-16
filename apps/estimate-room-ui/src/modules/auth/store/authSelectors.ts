import { AuthStates } from '../types';
import type { AuthState } from '../types';
import { authStateKey } from './authStore';

type AuthStateRoot = {
  readonly [authStateKey]: AuthState;
};

export const selectAuthState = (state: AuthStateRoot) => state[authStateKey];
export const selectAuthStatus = (state: AuthStateRoot) => selectAuthState(state).status;
export const selectAuthUser = (state: AuthStateRoot) => selectAuthState(state).user;
export const selectIsAuthenticated = (state: AuthStateRoot) =>
  selectAuthStatus(state) === AuthStates.AUTHENTICATED && selectAuthUser(state) !== null;
