import ArrowBackRoundedIcon from '@mui/icons-material/ArrowBackRounded';
import MailOutlineRoundedIcon from '@mui/icons-material/MailOutlineRounded';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import {
  AppAlert,
  AppBox,
  AppButton,
  AppLink,
  AppProgress,
  AppStack,
  AppTextField,
  AppTypography,
  OverlineText
} from '@/shared/ui';

import { createEmailValidationRules } from './utils';
import { AuthCard, AuthIntro, AuthShell } from './components';
import { useForgotPasswordPage } from './hooks';

export const ForgotPasswordPage = () => {
  const {
    form: {
      formState: { errors, isSubmitting, isValid },
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
          <AppStack spacing={3}>
            <AppTypography align="center" color="text.secondary" variant="body2">
              Please check your inbox and follow the instructions to reset your password.
            </AppTypography>
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
            <AppBox
              sx={{
                borderTop: (theme) => `1px solid ${theme.app.borders.ghost}`,
                pt: 3
              }}
            >
              <AppStack alignItems="center" spacing={1.5}>
                <AppLink
                  color="primary"
                  component={RouterLink}
                  sx={{ alignItems: 'center', display: 'inline-flex', gap: 1 }}
                  to={appRoutes.login}
                  underline="none"
                  variant="overline"
                >
                  <ArrowBackRoundedIcon fontSize="inherit" />
                  Back to Sign In
                </AppLink>
              </AppStack>
            </AppBox>
          </AppStack>
        ) : (
          <AppBox component="form" noValidate onSubmit={onSubmit}>
            <AppStack spacing={2.5}>
              {errors.root?.message ? <AppAlert severity="error">{errors.root.message}</AppAlert> : null}
              <AppStack spacing={1}>
                <OverlineText>Email Address</OverlineText>
                <AppTextField
                  autoComplete="email"
                  error={Boolean(errors.email)}
                  helperText={errors.email?.message}
                  placeholder="name@company.com"
                  type="email"
                  {...register('email', createEmailValidationRules())}
                />
              </AppStack>
              <AppButton
                disabled={!isValid}
                fullWidth
                loading={isSubmitting}
                loadingText="Sending Link..."
                type="submit"
                variant="contained"
              >
                Send Reset Link
              </AppButton>
              <AppLink
                color="primary"
                component={RouterLink}
                sx={{ alignItems: 'center', display: 'inline-flex', gap: 1, mx: 'auto' }}
                to={appRoutes.login}
                underline="none"
                variant="overline"
              >
                <ArrowBackRoundedIcon fontSize="inherit" />
                Back to Sign In
              </AppLink>
            </AppStack>
          </AppBox>
        )}
      </AuthCard>
      {isSubmitted ? (
        <AppTypography sx={{ mt: 3, textAlign: 'center' }} variant="body2">
          Didn&apos;t receive the email?{' '}
          <AppLink
            color="primary"
            component="button"
            disabled={isResending}
            onClick={() => {
              void onResend();
            }}
            sx={{
              alignItems: 'center',
              background: 'none',
              border: 0,
              cursor: 'pointer',
              display: 'inline-flex',
              gap: 0.75,
              p: 0
            }}
            type="button"
            underline="none"
          >
            {isResending ? <AppProgress color="inherit" size={14} /> : null}
            {isResending ? 'Resending...' : 'Click to resend link'}
          </AppLink>
        </AppTypography>
      ) : null}
    </AuthShell>
  );
};
