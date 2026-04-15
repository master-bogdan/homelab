import type { ThemeMode } from '@/shared/theme';

export interface SettingsSummary {
  readonly apiBaseUrl: string;
  readonly environment: string;
  readonly themeMode: ThemeMode;
  readonly wsBaseUrl: string;
}
