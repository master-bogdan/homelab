import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppButton, AppLink, AppStack, AppTypography, OverlineText } from '@/shared/ui';

import { AuthCard, AuthShell } from './components';

export const ResetPasswordSuccessPage = () => (
  <AuthShell>
    <AuthCard sx={{ mx: 'auto', maxWidth: 440 }}>
      <AppStack spacing={3} textAlign="center">
        <AppStack spacing={1.5}>
          <AppTypography component="h1" variant="h4">
            Password Updated
          </AppTypography>
          <AppTypography color="text.secondary" variant="body2">
            Your password has been successfully reset. You can now sign in with your
            new credentials.
          </AppTypography>
        </AppStack>
        <AppButton component={RouterLink} fullWidth to={appRoutes.login} variant="contained">
          Sign In
        </AppButton>
        <AppStack spacing={1}>
          <OverlineText>Need technical help?</OverlineText>
          <AppLink color="primary" href="#" onClick={(event) => event.preventDefault()} underline="none">
            Contact Architect Support
          </AppLink>
        </AppStack>
      </AppStack>
    </AuthCard>
  </AuthShell>
);
