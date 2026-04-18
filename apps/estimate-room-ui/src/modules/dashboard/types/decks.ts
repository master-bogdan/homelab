import type { DashboardDeckPresetKey } from './status';

export interface DashboardDeckPreset {
  readonly description: string;
  readonly key: DashboardDeckPresetKey;
  readonly label: string;
  readonly deck: {
    readonly kind: string;
    readonly name: string;
    readonly values: string[];
  };
}
