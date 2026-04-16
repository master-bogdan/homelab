import { dashboardReducer } from './dashboardSlice';

export const DASHBOARD_STATE_KEY = 'dashboard';

export const dashboardStore = {
  reducer: dashboardReducer,
  stateKey: DASHBOARD_STATE_KEY
} as const;
