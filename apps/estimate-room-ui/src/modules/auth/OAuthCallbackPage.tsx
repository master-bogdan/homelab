import ArrowBackRoundedIcon from '@mui/icons-material/ArrowBackRounded';
import { CircularProgress, Link, Stack, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppButton } from '@/shared/ui';

import { AuthCard, AuthShell } from './components';
import { useOAuthCallbackPage } from './hooks';

export const OAuthCallbackPage = () => {
  const { errorMessage, isLoading } = useOAuthCallbackPage();

  return (
    <AuthShell>
      <AuthCard sx={{ mx: 'auto', maxWidth: 440 }}>
        <Stack alignItems="center" spacing={3} textAlign="center">
          {isLoading ? <CircularProgress size={28} /> : null}
          <Stack spacing={1.5}>
            <Typography component="h1" variant="h4">
              {isLoading ? 'Signing You In' : 'Unable to Complete Sign-In'}
            </Typography>
            <Typography color="text.secondary" variant="body2">
              {isLoading
                ? 'Finalizing your EstimateRoom session.'
                : errorMessage ?? 'Please return to sign in and try again.'}
            </Typography>
          </Stack>
          {!isLoading ? (
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
          ) : null}
        </Stack>
      </AuthCard>
    </AuthShell>
  );
};
