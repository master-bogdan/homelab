import { useParams } from 'react-router-dom';

import { usePageTitle } from '@/shared/hooks';

import { teamsService } from '../api/teamsApi';

export const useTeamDetailsPage = () => {
  const { id = '' } = useParams();
  const team = teamsService.getTeamDetails(id);

  usePageTitle(team?.name ?? 'Team Details');

  return {
    team,
    teamId: id
  };
};
