import { combineReducers } from '@reduxjs/toolkit';

import { authReducer } from '@/modules/auth/store';

import { uiReducer } from './uiSlice';

export const rootReducer = combineReducers({
  auth: authReducer,
  ui: uiReducer
});
