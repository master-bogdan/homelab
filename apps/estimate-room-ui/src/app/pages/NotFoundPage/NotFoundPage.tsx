import { Link as RouterLink } from 'react-router-dom';

import { AppButton, AppPageState, AppStack, SectionCard } from '@/shared/components';
import { AppRoutes } from '@/app/router/routePaths';
import { usePageTitle } from '@/shared/hooks';

import { notFoundPageCardSx, notFoundPageRootSx } from './NotFoundPage.styles';

export const NotFoundPage = () => {
  usePageTitle('Not Found');

  return (
    <AppStack sx={notFoundPageRootSx}>
      <SectionCard
        description="The route does not exist in the current application scaffold."
        sx={notFoundPageCardSx}
        title="Page not found"
      >
        <AppPageState
          action={
            <AppButton component={RouterLink} to={AppRoutes.DASHBOARD} variant="contained">
              Go to dashboard
            </AppButton>
          }
          description="Double-check the path or head back to the dashboard entry point."
          title="Page not found"
          titleComponent="h2"
          titleVariant="body1"
        />
      </SectionCard>
    </AppStack>
  );
};
