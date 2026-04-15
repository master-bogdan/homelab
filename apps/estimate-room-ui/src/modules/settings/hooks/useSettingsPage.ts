import { useAppDispatch, useAppSelector } from '@/shared/store';
import { toggleThemeMode } from '@/modules/system';
import { selectThemeMode } from '@/modules/system';
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
