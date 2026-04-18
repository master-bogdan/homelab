import type { UnknownAction } from '@reduxjs/toolkit';
import type { ThunkDispatch } from 'redux-thunk';

import type { rootReducer } from '@/app/store/rootReducer';

export type RootState = ReturnType<typeof rootReducer>;
export type AppDispatch = ThunkDispatch<RootState, undefined, UnknownAction>;
