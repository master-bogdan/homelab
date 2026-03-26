import { Stack, Typography } from '@mui/material';

import { AppPageState, SectionCard } from '@/shared/ui';

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
        <AppPageState
          description="Connect the team details API to show membership, access, and activity here."
          title="No team data yet"
          titleVariant="body1"
        />
      )}
    </SectionCard>
  );
};
