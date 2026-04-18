import type { ThemeMode } from '@/shared/types';

export interface SettingsSummary {
  readonly apiBaseUrl: string;
  readonly environment: string;
  readonly themeMode: ThemeMode;
  readonly wsBaseUrl: string;
}
