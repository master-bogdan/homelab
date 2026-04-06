import { Stack } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { AppButton, AppPageState, SectionCard } from '@/shared/ui';
import { appRoutes } from '@/shared/constants/routes';
import { usePageTitle } from '@/shared/hooks';

export const NotFoundPage = () => {
  usePageTitle('Not Found');

  return (
    <Stack
      sx={{
        display: 'grid',
        placeItems: 'center',
        minHeight: '100vh',
        px: 3,
        py: 6
      }}
    >
      <SectionCard
        description="The route does not exist in the current application scaffold."
        sx={{ maxWidth: 560, width: '100%' }}
        title="Page not found"
      >
        <AppPageState
          action={
            <AppButton component={RouterLink} to={appRoutes.dashboard} variant="contained">
              Go to dashboard
            </AppButton>
          }
          description="Double-check the path or head back to the dashboard entry point."
          title="Page not found"
          titleComponent="h2"
          titleVariant="body1"
        />
      </SectionCard>
    </Stack>
  );
};
