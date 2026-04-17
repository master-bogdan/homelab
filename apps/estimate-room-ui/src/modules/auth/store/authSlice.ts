import { createSlice, type PayloadAction } from '@reduxjs/toolkit';

import { apiSessionExpired } from '@/shared/api/sessionLifecycle';
import type { AuthUser } from '@/shared/types';

import { AuthStates } from '../constants';
import type { AuthState } from '../types';

const initialState: AuthState = {
  oauthCallback: {
    errorMessage: null,
    redirectTo: null,
    requestKey: null,
    status: 'idle'
  },
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
    setOAuthCallbackFailed: (state, action: PayloadAction<string>) => {
      state.oauthCallback.errorMessage = action.payload;
      state.oauthCallback.redirectTo = null;
      state.oauthCallback.status = 'failed';
    },
    setOAuthCallbackPending: (state, action: PayloadAction<string>) => {
      state.oauthCallback.errorMessage = null;
      state.oauthCallback.redirectTo = null;
      state.oauthCallback.requestKey = action.payload;
      state.oauthCallback.status = 'pending';
    },
    setOAuthCallbackSucceeded: (
      state,
      action: PayloadAction<{ readonly redirectTo: string; readonly requestKey: string }>
    ) => {
      state.oauthCallback.errorMessage = null;
      state.oauthCallback.redirectTo = action.payload.redirectTo;
      state.oauthCallback.requestKey = action.payload.requestKey;
      state.oauthCallback.status = 'succeeded';
    },
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

export const {
  clearSession,
  hydrateSession,
  setOAuthCallbackFailed,
  setOAuthCallbackPending,
  setOAuthCallbackSucceeded,
  setSession
} = authSlice.actions;
export const authReducer = authSlice.reducer;
