import ArrowBackRoundedIcon from '@mui/icons-material/ArrowBackRounded';
import { Alert, Box, Link, Stack, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppButton, AppPageState, OverlineText } from '@/shared/ui';

import { createPasswordValidationRules } from './utils';
import {
  AuthCard,
  AuthIntro,
  AuthShell,
  PasswordField,
  PasswordRecommendations
} from './components';
import { useResetPasswordPage } from './hooks';

export const ResetPasswordPage = () => {
  const {
    form: {
      formState: { errors, isSubmitting, isValid },
      register
    },
    invalidLinkCopy,
    onSubmit,
    pageError,
    pageState,
    password
  } = useResetPasswordPage();

  if (pageState === 'invalid') {
    return (
      <AuthShell pattern="dots">
        <AuthCard sx={{ mx: 'auto', maxWidth: 440 }}>
          <Stack spacing={3} textAlign="center">
            <Stack spacing={1.5}>
              <Typography component="h1" variant="h4">
                {invalidLinkCopy.title}
              </Typography>
              <Typography color="text.secondary" variant="body2">
                {pageError ?? invalidLinkCopy.description}
              </Typography>
            </Stack>
            <AppButton component={RouterLink} fullWidth to={appRoutes.forgotPassword} variant="contained">
              Request New Link
            </AppButton>
            <Link
              color="text.secondary"
              component={RouterLink}
              sx={{ alignItems: 'center', display: 'inline-flex', gap: 1, justifyContent: 'center' }}
              to={appRoutes.login}
              underline="none"
              variant="body2"
            >
              <ArrowBackRoundedIcon fontSize="inherit" />
              Back to Sign In
            </Link>
          </Stack>
        </AuthCard>
      </AuthShell>
    );
  }

  return (
    <AuthShell>
      <AuthIntro
        description="Create a strong new password for your EstimateRoom account."
        title={pageState === 'validating' ? 'Validating Link' : 'Set New Password'}
      />
      <AuthCard>
        {pageState === 'validating' ? (
          <AppPageState
            description="Confirming your password reset link before you choose a new password."
            isLoading
            title="Validating Link"
          />
        ) : (
          <Box component="form" noValidate onSubmit={onSubmit}>
            <Stack spacing={2.5}>
              {errors.root?.message ? <Alert severity="error">{errors.root.message}</Alert> : null}
              <Stack spacing={1}>
                <OverlineText>New Password</OverlineText>
                <PasswordField
                  autoComplete="new-password"
                  error={Boolean(errors.password)}
                  fullWidth
                  helperText={errors.password?.message}
                  placeholder="••••••••"
                  {...register('password', createPasswordValidationRules())}
                />
              </Stack>
              <Stack spacing={1}>
                <OverlineText>Confirm New Password</OverlineText>
                <PasswordField
                  autoComplete="new-password"
                  error={Boolean(errors.confirmPassword)}
                  fullWidth
                  helperText={errors.confirmPassword?.message}
                  placeholder="••••••••"
                  {...register('confirmPassword', {
                    required: 'Please confirm your new password.',
                    validate: (value, values) =>
                      value === values.password || 'Passwords do not match.'
                  })}
                />
              </Stack>
              <PasswordRecommendations password={password} />
              <AppButton
                disabled={!isValid}
                fullWidth
                loading={isSubmitting}
                loadingText="Resetting Password..."
                type="submit"
                variant="contained"
              >
                Reset Password
              </AppButton>
            </Stack>
          </Box>
        )}
      </AuthCard>
      <Typography sx={{ mt: 3, textAlign: 'center' }} variant="body2">
        Remember your password?{' '}
        <Link color="primary" component={RouterLink} to={appRoutes.login} underline="none">
          Back to Login
        </Link>
      </Typography>
    </AuthShell>
  );
};
