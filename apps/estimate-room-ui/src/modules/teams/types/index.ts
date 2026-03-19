import type { Team } from '@/shared/types';

export interface TeamPageData {
  readonly team: Team | null;
  readonly teamId: string;
}
