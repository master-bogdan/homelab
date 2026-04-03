import type { RootState } from '@/app/store/store';

export const selectAuthState = (state: RootState) => state.auth;
export const selectAuthStatus = (state: RootState) => state.auth.status;
export const selectAuthUser = (state: RootState) => state.auth.user;
export const selectIsAuthenticated = (state: RootState) =>
  state.auth.status === 'authenticated' && state.auth.user !== null;
