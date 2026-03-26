import GitHubIcon from '@mui/icons-material/GitHub';
import { Alert, Box, Link, Stack, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppButton, AppTextField, OverlineText } from '@/shared/ui';

import { createEmailValidationRules } from './utils';
import { AuthActionDivider, AuthCard, AuthIntro, AuthShell, PasswordField } from './components';
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
    <AuthShell>
      <AuthIntro
        description="Access your precision estimation workspace."
        title="Welcome Back"
      />
      <AuthCard>
        <Box component="form" noValidate onSubmit={onSubmit}>
          <Stack spacing={2.5}>
            {errors.root?.message ? <Alert severity="error">{errors.root.message}</Alert> : null}
            <Stack spacing={1}>
              <OverlineText>Email Address</OverlineText>
              <AppTextField
                autoComplete="email"
                error={Boolean(errors.email)}
                helperText={errors.email?.message}
                placeholder="name@company.com"
                type="email"
                {...register('email', createEmailValidationRules())}
              />
            </Stack>

            <Stack spacing={1}>
              <Stack alignItems="center" direction="row" justifyContent="space-between">
                <OverlineText>Password</OverlineText>
                <Link
                  color="primary"
                  component={RouterLink}
                  to={appRoutes.forgotPassword}
                  underline="none"
                  variant="overline"
                >
                  Forgot?
                </Link>
              </Stack>
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
            </Stack>

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
          </Stack>
        </Box>
      </AuthCard>
      <Typography sx={{ mt: 3, textAlign: 'center' }} variant="body2">
        Don&apos;t have an account?{' '}
        <Link color="primary" component={RouterLink} to={appRoutes.register} underline="none">
          Register now
        </Link>
      </Typography>
    </AuthShell>
  );
};
