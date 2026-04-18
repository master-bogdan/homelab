import { useTeamDetailsPage } from '@/modules/teams/hooks/useTeamDetailsPage';
import { AppPageState, AppStack, AppTypography, SectionCard } from '@/shared/components';

export const TeamDetailsPage = () => {
  const { team, teamId } = useTeamDetailsPage();

  return (
    <SectionCard
      description="Future team-level membership, permissions, and room ownership data can live here."
      title={`Team ${teamId}`}
    >
      {team ? (
        <AppStack spacing={1}>
          <AppTypography variant="h4">{team.name}</AppTypography>
          <AppTypography color="text.secondary" variant="body2">
            Region: {team.region}
          </AppTypography>
          <AppTypography color="text.secondary" variant="body2">
            Members: {team.memberCount}
          </AppTypography>
        </AppStack>
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
