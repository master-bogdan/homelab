import type { ThemeMode } from '@/shared/theme';

import type { SystemUiState } from '../types';
import { systemStateKey } from './system.store';

type SystemUiStateRoot = {
  readonly [systemStateKey]: {
    readonly ui: SystemUiState;
  };
};

export const selectSystemUiState = (state: SystemUiStateRoot) =>
  state[systemStateKey].ui;

export const selectIsSidebarOpen = (state: SystemUiStateRoot) =>
  selectSystemUiState(state).sidebarOpen;

export const selectThemeMode = (state: SystemUiStateRoot): ThemeMode =>
  selectSystemUiState(state).themeMode;
