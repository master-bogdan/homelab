import { createSlice, type PayloadAction } from '@reduxjs/toolkit';

import type { AuthUser } from '@/shared/types';

import { AUTH_STATUSES, type AuthState } from '../types';

const initialState: AuthState = {
  status: AUTH_STATUSES.UNKNOWN,
  user: null
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearSession: (state) => {
      state.status = AUTH_STATUSES.UNAUTHENTICATED;
      state.user = null;
    },
    hydrateSession: (_state, action: PayloadAction<AuthState>) => action.payload,
    setSession: (state, action: PayloadAction<AuthUser>) => {
      state.status = AUTH_STATUSES.AUTHENTICATED;
      state.user = action.payload;
    }
  }
});

export const { clearSession, hydrateSession, setSession } = authSlice.actions;
export const authReducer = authSlice.reducer;
