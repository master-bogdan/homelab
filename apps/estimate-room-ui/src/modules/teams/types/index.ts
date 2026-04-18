import type { Team } from './team';

export interface TeamPageData {
  readonly team: Team | null;
  readonly teamId: string;
}

export type { Team } from './team';
