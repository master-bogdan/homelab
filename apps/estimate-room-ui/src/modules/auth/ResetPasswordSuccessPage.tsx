import { Link, Stack, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppButton, OverlineText } from '@/shared/ui';

import { AuthCard, AuthShell } from './components';

export const ResetPasswordSuccessPage = () => (
  <AuthShell>
    <AuthCard sx={{ mx: 'auto', maxWidth: 440 }}>
      <Stack spacing={3} textAlign="center">
        <Stack spacing={1.5}>
          <Typography component="h1" variant="h4">
            Password Updated
          </Typography>
          <Typography color="text.secondary" variant="body2">
            Your password has been successfully reset. You can now sign in with your
            new credentials.
          </Typography>
        </Stack>
        <AppButton component={RouterLink} fullWidth to={appRoutes.login} variant="contained">
          Sign In
        </AppButton>
        <Stack spacing={1}>
          <OverlineText>Need technical help?</OverlineText>
          <Link color="primary" href="#" onClick={(event) => event.preventDefault()} underline="none">
            Contact Architect Support
          </Link>
        </Stack>
      </Stack>
    </AuthCard>
  </AuthShell>
);
