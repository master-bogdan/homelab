import { configureStore } from '@reduxjs/toolkit';

import { rtkQueryMiddleware } from './middleware';
import { rootReducer } from './rootReducer';

export const store = configureStore({
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(rtkQueryMiddleware),
  reducer: rootReducer,
  devTools: import.meta.env.DEV
});
