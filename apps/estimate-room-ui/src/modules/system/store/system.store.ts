import { systemReducer } from './systemSlice';

export const systemStateKey = 'system';

export const systemStore = {
  reducer: systemReducer,
  stateKey: systemStateKey
} as const;
