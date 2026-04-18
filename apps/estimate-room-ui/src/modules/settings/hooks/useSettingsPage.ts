import { useAppDispatch, useAppSelector } from '@/shared/hooks';
import { toggleThemeMode } from '@/modules/system';
import { selectThemeMode } from '@/modules/system';
import { AppConfig } from '@/config';
import { usePageTitle } from '@/shared/hooks';

export const useSettingsPage = () => {
  const dispatch = useAppDispatch();
  const themeMode = useAppSelector(selectThemeMode);

  usePageTitle('Settings');

  return {
    apiBaseUrl: AppConfig.API_BASE_URL,
    environment: AppConfig.ENVIRONMENT,
    themeMode,
    toggleTheme: () => dispatch(toggleThemeMode()),
    wsBaseUrl: AppConfig.WS_BASE_URL
  };
};
