import { createSlice, type PayloadAction } from '@reduxjs/toolkit';

import type { AuthUser } from '@/shared/types';

import type { AuthState } from '../types';

const initialState: AuthState = {
  status: 'unknown',
  user: null
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearSession: (state) => {
      state.status = 'unauthenticated';
      state.user = null;
    },
    hydrateSession: (_state, action: PayloadAction<AuthState>) => action.payload,
    setSession: (state, action: PayloadAction<AuthUser>) => {
      state.status = 'authenticated';
      state.user = action.payload;
    }
  }
});

export const { clearSession, hydrateSession, setSession } = authSlice.actions;
export const authReducer = authSlice.reducer;
