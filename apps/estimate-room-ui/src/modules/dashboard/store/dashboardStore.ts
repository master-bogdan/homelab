import { dashboardReducer } from './dashboardSlice';

export const dashboardStateKey = 'dashboard';

export const dashboardStore = {
  reducer: dashboardReducer,
  stateKey: dashboardStateKey
} as const;
