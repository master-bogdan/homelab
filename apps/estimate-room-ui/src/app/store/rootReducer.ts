import { combineReducers } from '@reduxjs/toolkit';

import { authReducer } from '@/modules/auth/store';
import { dashboardReducer } from '@/modules/dashboard/store';
import { systemReducer } from '@/modules/system/store';
import { api } from '@/shared/api';

import { uiReducer } from './uiSlice';

export const rootReducer = combineReducers({
  [api.reducerPath]: api.reducer,
  auth: authReducer,
  dashboard: dashboardReducer,
  system: systemReducer,
  ui: uiReducer
});
