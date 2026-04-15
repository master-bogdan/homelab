import { combineReducers } from '@reduxjs/toolkit';

import { authStore } from '@/modules/auth/store/authStore';
import { dashboardStore } from '@/modules/dashboard/store/dashboardStore';
import { systemStore } from '@/modules/system/store/systemStore';
import { api } from '@/shared/api';

export const rootReducer = combineReducers({
  [api.reducerPath]: api.reducer,
  [authStore.stateKey]: authStore.reducer,
  [dashboardStore.stateKey]: dashboardStore.reducer,
  [systemStore.stateKey]: systemStore.reducer
});
