import { createSlice, type PayloadAction } from '@reduxjs/toolkit';

import type { ThemeMode } from '@/theme';

export interface UiState {
  readonly sidebarOpen: boolean;
  readonly themeMode: ThemeMode;
}

const initialState: UiState = {
  sidebarOpen: true,
  themeMode: 'light'
};

const uiSlice = createSlice({
  name: 'ui',
  initialState,
  reducers: {
    closeSidebar: (state) => {
      state.sidebarOpen = false;
    },
    openSidebar: (state) => {
      state.sidebarOpen = true;
    },
    setSidebarOpen: (state, action: PayloadAction<boolean>) => {
      state.sidebarOpen = action.payload;
    },
    setThemeMode: (state, action: PayloadAction<ThemeMode>) => {
      state.themeMode = action.payload;
    },
    toggleThemeMode: (state) => {
      state.themeMode = state.themeMode === 'light' ? 'dark' : 'light';
    }
  }
});

export const {
  closeSidebar,
  openSidebar,
  setSidebarOpen,
  setThemeMode,
  toggleThemeMode
} = uiSlice.actions;

export const uiReducer = uiSlice.reducer;
