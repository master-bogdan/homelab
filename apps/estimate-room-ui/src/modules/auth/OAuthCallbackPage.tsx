import ArrowBackRoundedIcon from '@mui/icons-material/ArrowBackRounded';
import { Link } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppButton, AppPageState } from '@/shared/ui';

import { AuthCard, AuthShell } from './components';
import { useOAuthCallbackPage } from './hooks';

export const OAuthCallbackPage = () => {
  const { errorMessage, isLoading } = useOAuthCallbackPage();

  return (
    <AuthShell>
      <AuthCard sx={{ mx: 'auto', maxWidth: 440 }}>
        <AppPageState
          action={
            !isLoading ? (
              <>
                <AppButton component={RouterLink} fullWidth to={appRoutes.login} variant="contained">
                  Return to Sign In
                </AppButton>
                <Link
                  color="text.secondary"
                  component={RouterLink}
                  sx={{
                    alignItems: 'center',
                    display: 'inline-flex',
                    gap: 1,
                    justifyContent: 'center'
                  }}
                  to={appRoutes.login}
                  underline="none"
                  variant="body2"
                >
                  <ArrowBackRoundedIcon fontSize="inherit" />
                  Back to Sign In
                </Link>
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
    </AuthShell>
  );
};
