import GitHubIcon from '@mui/icons-material/GitHub';
import { Alert, Box, Link, Stack, TextField, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppButton, OverlineText } from '@/shared/ui';

import { AuthActionDivider, AuthCard, AuthIntro, AuthShell } from './components';
import { useRegisterPage } from './hooks';

export const RegisterPage = () => {
  const {
    form: {
      formState: { errors, isSubmitting },
      register
    },
    onSubmit,
    onSubmitWithGithub
  } = useRegisterPage();

  return (
    <AuthShell>
      <AuthIntro
        description="Define your professional profile to start estimating."
        title="Create Workspace"
      />
      <AuthCard>
        <Box component="form" noValidate onSubmit={onSubmit}>
          <Stack spacing={2.5}>
            {errors.root?.message ? <Alert severity="error">{errors.root.message}</Alert> : null}

            <Stack spacing={1}>
              <OverlineText>Full Name</OverlineText>
              <TextField
                autoComplete="name"
                error={Boolean(errors.displayName)}
                fullWidth
                helperText={errors.displayName?.message}
                placeholder="Jane Architect"
                {...register('displayName', {
                  maxLength: {
                    message: 'Display name must be 100 characters or less.',
                    value: 100
                  },
                  required: 'Full name is required.'
                })}
              />
            </Stack>

            <Stack spacing={1}>
              <OverlineText>Work Email</OverlineText>
              <TextField
                autoComplete="email"
                error={Boolean(errors.email)}
                fullWidth
                helperText={errors.email?.message}
                placeholder="name@company.com"
                {...register('email', {
                  pattern: {
                    message: 'Enter a valid email address.',
                    value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/u
                  },
                  required: 'Email is required.'
                })}
              />
            </Stack>

            <Stack spacing={1}>
              <OverlineText>Password</OverlineText>
              <TextField
                autoComplete="new-password"
                error={Boolean(errors.password)}
                fullWidth
                helperText={errors.password?.message}
                placeholder="••••••••"
                type="password"
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
              fullWidth
              loading={isSubmitting}
              loadingText="Creating Account..."
              type="submit"
              variant="contained"
            >
              Initialize Account
            </AppButton>

            <AuthActionDivider />

            <AppButton
              color="secondary"
              fullWidth
              loading={isSubmitting}
              onClick={() => {
                void onSubmitWithGithub();
              }}
              startIcon={<GitHubIcon />}
              type="button"
              variant="contained"
            >
              Sign up with GitHub
            </AppButton>
          </Stack>
        </Box>
      </AuthCard>
      <Typography sx={{ mt: 3, textAlign: 'center' }} variant="body2">
        Already have an account?{' '}
        <Link color="primary" component={RouterLink} to={appRoutes.login} underline="none">
          Sign In
        </Link>
      </Typography>
    </AuthShell>
  );
};
