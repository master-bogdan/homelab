import { combineReducers } from '@reduxjs/toolkit';

import { authReducer } from '@/modules/auth/store/slice';
import { AUTH_STATE_KEY } from '@/modules/auth/store/types';
import { dashboardReducer } from '@/modules/dashboard/store/slice';
import { DASHBOARD_STATE_KEY } from '@/modules/dashboard/store/types';
import { systemReducer } from '@/modules/system/store/slice';
import { systemStateKey } from '@/modules/system/store/types';
import { api } from '@/shared/api';

export const rootReducer = combineReducers({
  [api.reducerPath]: api.reducer,
  [AUTH_STATE_KEY]: authReducer,
  [DASHBOARD_STATE_KEY]: dashboardReducer,
  [systemStateKey]: systemReducer
});
