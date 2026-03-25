import ArrowBackRoundedIcon from '@mui/icons-material/ArrowBackRounded';
import MailOutlineRoundedIcon from '@mui/icons-material/MailOutlineRounded';
import { Alert, Box, Link, Stack, TextField, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppButton, OverlineText } from '@/shared/ui';

import { createEmailValidationRules } from './utils';
import { AuthCard, AuthIntro, AuthShell } from './components';
import { useForgotPasswordPage } from './hooks';

export const ForgotPasswordPage = () => {
  const {
    form: {
      formState: { errors, isSubmitting },
      register
    },
    isResending,
    isSubmitted,
    onResend,
    onSubmit
  } = useForgotPasswordPage();

  return (
    <AuthShell>
      <AuthIntro
        description={
          isSubmitted
            ? "We've sent a password reset link to your email address."
            : 'Enter your work email and we will send you a secure reset link.'
        }
        title={isSubmitted ? 'Check Your Email' : 'Reset Password'}
      />
      <AuthCard>
        {isSubmitted ? (
          <Stack spacing={3}>
            <Typography align="center" color="text.secondary" variant="body2">
              Please check your inbox and follow the instructions to reset your password.
            </Typography>
            <AppButton
              color="secondary"
              component="a"
              fullWidth
              href="mailto:"
              startIcon={<MailOutlineRoundedIcon />}
              variant="contained"
            >
              Open Email App
            </AppButton>
            <Box
              sx={{
                borderTop: (theme) => `1px solid ${theme.app.borders.ghost}`,
                pt: 3
              }}
            >
              <Stack alignItems="center" spacing={1.5}>
                <Link
                  color="primary"
                  component={RouterLink}
                  sx={{ alignItems: 'center', display: 'inline-flex', gap: 1 }}
                  to={appRoutes.login}
                  underline="none"
                  variant="overline"
                >
                  <ArrowBackRoundedIcon fontSize="inherit" />
                  Back to Sign In
                </Link>
              </Stack>
            </Box>
          </Stack>
        ) : (
          <Box component="form" noValidate onSubmit={onSubmit}>
            <Stack spacing={2.5}>
              {errors.root?.message ? <Alert severity="error">{errors.root.message}</Alert> : null}
              <Stack spacing={1}>
                <OverlineText>Email Address</OverlineText>
                <TextField
                  autoComplete="email"
                  error={Boolean(errors.email)}
                  fullWidth
                  helperText={errors.email?.message}
                  placeholder="name@company.com"
                  type="email"
                  {...register('email', createEmailValidationRules())}
                />
              </Stack>
              <AppButton
                fullWidth
                loading={isSubmitting}
                loadingText="Sending Link..."
                type="submit"
                variant="contained"
              >
                Send Reset Link
              </AppButton>
              <Link
                color="primary"
                component={RouterLink}
                sx={{ alignItems: 'center', display: 'inline-flex', gap: 1, mx: 'auto' }}
                to={appRoutes.login}
                underline="none"
                variant="overline"
              >
                <ArrowBackRoundedIcon fontSize="inherit" />
                Back to Sign In
              </Link>
            </Stack>
          </Box>
        )}
      </AuthCard>
      {isSubmitted ? (
        <Typography sx={{ mt: 3, textAlign: 'center' }} variant="body2">
          Didn&apos;t receive the email?{' '}
          <Link
            color="primary"
            component="button"
            disabled={isResending}
            onClick={() => {
              void onResend();
            }}
            sx={{ background: 'none', border: 0, cursor: 'pointer', p: 0 }}
            type="button"
            underline="none"
          >
            {isResending ? 'Resending...' : 'Click to resend link'}
          </Link>
        </Typography>
      ) : null}
    </AuthShell>
  );
};
