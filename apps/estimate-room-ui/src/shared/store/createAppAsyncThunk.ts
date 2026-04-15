import { createAsyncThunk } from '@reduxjs/toolkit';

import type { AppDispatch, RootState } from './types';

export const createAppAsyncThunk = createAsyncThunk.withTypes<{
  dispatch: AppDispatch;
  rejectValue: string;
  state: RootState;
}>();
