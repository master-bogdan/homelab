import { Stack, Typography } from '@mui/material';

import { SectionCard } from '@/shared/ui';

import { useTeamDetailsPage } from './hooks/useTeamDetailsPage';

export const TeamDetailsPage = () => {
  const { team, teamId } = useTeamDetailsPage();

  return (
    <SectionCard
      description="Future team-level membership, permissions, and room ownership data can live here."
      title={`Team ${teamId}`}
    >
      {team ? (
        <Stack spacing={1}>
          <Typography variant="h4">{team.name}</Typography>
          <Typography color="text.secondary" variant="body2">
            Region: {team.region}
          </Typography>
          <Typography color="text.secondary" variant="body2">
            Members: {team.memberCount}
          </Typography>
        </Stack>
      ) : (
        <Typography color="text.secondary" variant="body2">
          No team data is available for this route yet.
        </Typography>
      )}
    </SectionCard>
  );
};
