import GitHubIcon from '@mui/icons-material/GitHub';
import { Link as RouterLink } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';
import {
  AppAlert,
  AppBox,
  AppButton,
  AppLink,
  AppStack,
  AppTextField,
  AppTypography,
  OverlineText
} from '@/shared/ui';

import { createEmailValidationRules } from './utils';
import { AuthActionDivider, AuthCard, AuthIntro, AuthPageLayout, PasswordField } from './components';
import { useLoginPage } from './hooks';

export const LoginPage = () => {
  const {
    form: {
      formState: { errors, isSubmitting, isValid },
      register
    },
    isGithubLoading,
    onSubmit,
    onSubmitWithGithub
  } = useLoginPage();

  return (
    <AuthPageLayout>
      <AuthIntro
        description="Access your precision estimation workspace."
        title="Welcome Back"
      />
      <AuthCard>
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

            <AppStack spacing={1}>
              <AppStack alignItems="center" direction="row" justifyContent="space-between">
                <OverlineText>Password</OverlineText>
                <AppLink
                  color="primary"
                  component={RouterLink}
                  to={AppRoutes.FORGOT_PASSWORD}
                  underline="none"
                  variant="overline"
                >
                  Forgot?
                </AppLink>
              </AppStack>
              <PasswordField
                autoComplete="current-password"
                error={Boolean(errors.password)}
                fullWidth
                helperText={errors.password?.message}
                placeholder="••••••••"
                {...register('password', {
                  minLength: {
                    message: 'Password must be at least 8 characters.',
                    value: 8
                  },
                  required: 'Password is required.'
                })}
              />
            </AppStack>

            <AppButton
              disabled={!isValid || isGithubLoading}
              fullWidth
              loading={isSubmitting}
              loadingText="Signing In..."
              type="submit"
              variant="contained"
            >
              Sign In
            </AppButton>

            <AuthActionDivider />

            <AppButton
              color="secondary"
              disabled={isSubmitting}
              fullWidth
              loading={isGithubLoading}
              loadingText="Redirecting to GitHub..."
              onClick={() => {
                void onSubmitWithGithub();
              }}
              startIcon={<GitHubIcon />}
              type="button"
              variant="contained"
            >
              Continue with GitHub
            </AppButton>
          </AppStack>
        </AppBox>
      </AuthCard>
      <AppTypography sx={{ mt: 3, textAlign: 'center' }} variant="body2">
        Don&apos;t have an account?{' '}
        <AppLink color="primary" component={RouterLink} to={AppRoutes.REGISTER} underline="none">
          Register now
        </AppLink>
      </AppTypography>
    </AuthPageLayout>
  );
};
