import { authReducer } from './authSlice';

export const authStateKey = 'auth';

export const authStore = {
  reducer: authReducer,
  stateKey: authStateKey
} as const;
