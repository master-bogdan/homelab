import ArrowBackRoundedIcon from '@mui/icons-material/ArrowBackRounded';
import { Link as RouterLink } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';
import { AppButton, AppLink, AppPageState } from '@/shared/ui';

import { AuthCard, AuthPageLayout } from './components';
import { useOAuthCallbackPage } from './hooks';

export const OAuthCallbackPage = () => {
  const { errorMessage, isLoading } = useOAuthCallbackPage();

  return (
    <AuthPageLayout>
      <AuthCard sx={{ mx: 'auto', maxWidth: 440 }}>
        <AppPageState
          action={
            !isLoading ? (
              <>
                <AppButton component={RouterLink} fullWidth to={AppRoutes.LOGIN} variant="contained">
                  Return to Sign In
                </AppButton>
                <AppLink
                  color="text.secondary"
                  component={RouterLink}
                  sx={{
                    alignItems: 'center',
                    display: 'inline-flex',
                    gap: 1,
                    justifyContent: 'center'
                  }}
                  to={AppRoutes.LOGIN}
                  underline="none"
                  variant="body2"
                >
                  <ArrowBackRoundedIcon fontSize="inherit" />
                  Back to Sign In
                </AppLink>
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
      </AuthCard>
    </AuthPageLayout>
  );
};
