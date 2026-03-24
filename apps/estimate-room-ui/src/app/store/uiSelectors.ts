import type { RootState } from './store';

export const selectIsSidebarOpen = (state: RootState) => state.ui.sidebarOpen;
export const selectThemeMode = (state: RootState) => state.ui.themeMode;
