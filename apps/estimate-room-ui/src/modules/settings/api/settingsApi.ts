import type { SettingsSummary } from '../types';

export const settingsService = {
  savePreferences: async (settings: SettingsSummary) => Promise.resolve(settings)
};
