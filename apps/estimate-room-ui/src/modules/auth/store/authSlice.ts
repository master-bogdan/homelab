import { createSlice, type PayloadAction } from '@reduxjs/toolkit';

import { apiSessionExpired } from '@/shared/api/sessionLifecycle';
import type { AuthUser } from '@/shared/types';

import { AuthStates, type AuthState } from '../types';

const initialState: AuthState = {
  status: AuthStates.UNKNOWN,
  user: null
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearSession: (state) => {
      state.status = AuthStates.UNAUTHENTICATED;
      state.user = null;
    },
    hydrateSession: (_state, action: PayloadAction<AuthState>) => action.payload,
    setSession: (state, action: PayloadAction<AuthUser>) => {
      state.status = AuthStates.AUTHENTICATED;
      state.user = action.payload;
    }
  },
  extraReducers: (builder) => {
    builder.addCase(apiSessionExpired, (state) => {
      state.status = AuthStates.UNAUTHENTICATED;
      state.user = null;
    });
  }
});

export const { clearSession, hydrateSession, setSession } = authSlice.actions;
export const authReducer = authSlice.reducer;
