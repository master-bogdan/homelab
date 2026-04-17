import { authReducer } from './authSlice';

export const AUTH_STATE_KEY = 'auth';

export const authStore = {
  reducer: authReducer,
  stateKey: AUTH_STATE_KEY
} as const;
