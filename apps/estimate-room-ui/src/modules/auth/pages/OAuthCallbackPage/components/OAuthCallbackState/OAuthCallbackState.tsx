import { Link as RouterLink } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';
import { AppButton, AppPageState } from '@/shared/ui';

import { AuthBackToSignInLink } from '../../../../components';

interface OAuthCallbackStateProps {
  readonly errorMessage: string | null;
  readonly isLoading: boolean;
}

export const OAuthCallbackState = ({
  errorMessage,
  isLoading
}: OAuthCallbackStateProps) => (
  <AppPageState
    action={
      !isLoading ? (
        <>
          <AppButton component={RouterLink} fullWidth to={AppRoutes.LOGIN} variant="contained">
            Return to Sign In
          </AppButton>
          <AuthBackToSignInLink color="text.secondary" placement="centered" variant="body2" />
        </>
      ) : null
    }
    description={
      isLoading
        ? 'Finalizing your EstimateRoom session.'
        : errorMessage ?? 'Please return to sign in and try again.'
    }
    isLoading={isLoading}
    title={isLoading ? 'Signing You In' : 'Unable to Complete Sign-In'}
    titleComponent="h1"
    titleVariant="h4"
  />
);
