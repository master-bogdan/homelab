import { useAppDispatch, useAppSelector } from '@/app/store/hooks';
import { toggleThemeMode } from '@/app/store/uiSlice';
import { selectThemeMode } from '@/app/store/uiSelectors';
import { appConfig } from '@/shared/config/env';
import { usePageTitle } from '@/shared/hooks';

export const useSettingsPage = () => {
  const dispatch = useAppDispatch();
  const themeMode = useAppSelector(selectThemeMode);

  usePageTitle('Settings');

  return {
    apiBaseUrl: appConfig.apiBaseUrl,
    environment: appConfig.environment,
    themeMode,
    toggleTheme: () => dispatch(toggleThemeMode()),
    wsBaseUrl: appConfig.wsBaseUrl
  };
};
