import ArrowBackRoundedIcon from '@mui/icons-material/ArrowBackRounded';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import {
  AppAlert,
  AppBox,
  AppButton,
  AppLink,
  AppPageState,
  AppStack,
  AppTypography,
  OverlineText
} from '@/shared/ui';

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
          <AppStack spacing={3} textAlign="center">
            <AppStack spacing={1.5}>
              <AppTypography component="h1" variant="h4">
                {invalidLinkCopy.title}
              </AppTypography>
              <AppTypography color="text.secondary" variant="body2">
                {pageError ?? invalidLinkCopy.description}
              </AppTypography>
            </AppStack>
            <AppButton component={RouterLink} fullWidth to={appRoutes.forgotPassword} variant="contained">
              Request New Link
            </AppButton>
            <AppLink
              color="text.secondary"
              component={RouterLink}
              sx={{ alignItems: 'center', display: 'inline-flex', gap: 1, justifyContent: 'center' }}
              to={appRoutes.login}
              underline="none"
              variant="body2"
            >
              <ArrowBackRoundedIcon fontSize="inherit" />
              Back to Sign In
            </AppLink>
          </AppStack>
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
          <AppBox component="form" noValidate onSubmit={onSubmit}>
            <AppStack spacing={2.5}>
              {errors.root?.message ? <AppAlert severity="error">{errors.root.message}</AppAlert> : null}
              <AppStack spacing={1}>
                <OverlineText>New Password</OverlineText>
                <PasswordField
                  autoComplete="new-password"
                  error={Boolean(errors.password)}
                  fullWidth
                  helperText={errors.password?.message}
                  placeholder="••••••••"
                  {...register('password', createPasswordValidationRules())}
                />
              </AppStack>
              <AppStack spacing={1}>
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
              </AppStack>
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
            </AppStack>
          </AppBox>
        )}
      </AuthCard>
      <AppTypography sx={{ mt: 3, textAlign: 'center' }} variant="body2">
        Remember your password?{' '}
        <AppLink color="primary" component={RouterLink} to={appRoutes.login} underline="none">
          Back to Login
        </AppLink>
      </AppTypography>
    </AuthShell>
  );
};
