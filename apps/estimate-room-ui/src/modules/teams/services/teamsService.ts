import type { Team } from '@/shared/types';

export const teamsService = {
  getTeamDetails: (teamId: string): Team | null =>
    teamId
      ? {
          id: teamId,
          memberCount: 11,
          name: 'Atlantic Operations',
          region: 'US East'
        }
      : null
};
