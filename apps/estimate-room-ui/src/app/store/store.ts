import { configureStore } from '@reduxjs/toolkit';

import { api } from '@/shared/api';

import { rootReducer } from './rootReducer';

export const store = configureStore({
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
  reducer: rootReducer,
  devTools: import.meta.env.DEV
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
