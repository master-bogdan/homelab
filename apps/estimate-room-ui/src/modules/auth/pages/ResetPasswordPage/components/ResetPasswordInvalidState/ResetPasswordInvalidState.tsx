import { Link as RouterLink } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';
import {
  AppButton,
  AppStack,
  AppTypography
} from '@/shared/ui';

import { AuthBackToSignInLink, AuthNarrowCard, AuthPageLayout } from '../../../../components';

interface ResetPasswordInvalidStateProps {
  readonly description: string;
  readonly pageError: string | null;
  readonly title: string;
}

export const ResetPasswordInvalidState = ({
  description,
  pageError,
  title
}: ResetPasswordInvalidStateProps) => (
  <AuthPageLayout pattern="dots">
    <AuthNarrowCard>
      <AppStack spacing={3} textAlign="center">
        <AppStack spacing={1.5}>
          <AppTypography component="h1" variant="h4">
            {title}
          </AppTypography>
          <AppTypography color="text.secondary" variant="body2">
            {pageError ?? description}
          </AppTypography>
        </AppStack>
        <AppButton component={RouterLink} fullWidth to={AppRoutes.FORGOT_PASSWORD} variant="contained">
          Request New Link
        </AppButton>
        <AuthBackToSignInLink color="text.secondary" placement="centered" variant="body2" />
      </AppStack>
    </AuthNarrowCard>
  </AuthPageLayout>
);
