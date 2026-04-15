import { AppStack, AppTypography, SectionCard } from '@/shared/ui';

import { useProfilePage } from './hooks/useProfilePage';

export const ProfilePage = () => {
  const { profile } = useProfilePage();

  return (
    <SectionCard
      description="Use this page for the authenticated user profile, access scope, and editable preferences."
      title="Profile"
    >
      <AppStack spacing={1}>
        <AppTypography variant="h4">{profile.displayName}</AppTypography>
        <AppTypography color="text.secondary" variant="body2">
          {profile.email}
        </AppTypography>
        <AppTypography color="text.secondary" variant="body2">
          Role: {profile.roleLabel}
        </AppTypography>
        <AppTypography color="text.secondary" variant="body2">
          Teams: {profile.teamCount}
        </AppTypography>
      </AppStack>
    </SectionCard>
  );
};
