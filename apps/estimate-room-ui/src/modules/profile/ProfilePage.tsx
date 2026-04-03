import { Stack, Typography } from '@mui/material';

import { SectionCard } from '@/shared/ui';

import { useProfilePage } from './hooks/useProfilePage';

export const ProfilePage = () => {
  const { profile } = useProfilePage();

  return (
    <SectionCard
      description="Use this page for the authenticated user profile, access scope, and editable preferences."
      title="Profile"
    >
      <Stack spacing={1}>
        <Typography variant="h4">{profile.displayName}</Typography>
        <Typography color="text.secondary" variant="body2">
          {profile.email}
        </Typography>
        <Typography color="text.secondary" variant="body2">
          Role: {profile.roleLabel}
        </Typography>
        <Typography color="text.secondary" variant="body2">
          Teams: {profile.teamCount}
        </Typography>
      </Stack>
    </SectionCard>
  );
};
