import ArrowForwardRoundedIcon from '@mui/icons-material/ArrowForwardRounded';
import GitHubIcon from '@mui/icons-material/GitHub';
import { Alert, Box, Link, Stack, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppButton, AppTextField, OverlineText } from '@/shared/ui';

import { createEmailValidationRules, createPasswordValidationRules } from './utils';
import {
  AuthActionDivider,
  AuthCard,
  AuthIntro,
  AuthShell,
  PasswordField,
  PasswordRecommendations
} from './components';
import { useRegisterPage } from './hooks';

export const RegisterPage = () => {
  const {
    form: {
      formState: { errors, isSubmitting, isValid },
      register
    },
    isGithubLoading,
    onSubmit,
    onSubmitWithGithub,
    password
  } = useRegisterPage();

  return (
    <AuthShell>
      <AuthIntro
        description="Define your professional profile to start estimating."
        title="Create Workspace"
      />
      <AuthCard sx={{ maxWidth: 460, mx: 'auto' }}>
        <Box component="form" noValidate onSubmit={onSubmit}>
          <Stack spacing={2.5}>
            {errors.root?.message ? <Alert severity="error">{errors.root.message}</Alert> : null}

            <Stack spacing={1}>
              <OverlineText>Full Name</OverlineText>
              <AppTextField
                autoComplete="name"
                error={Boolean(errors.displayName)}
                helperText={errors.displayName?.message}
                placeholder="John Doe"
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
              <AppTextField
                autoComplete="email"
                error={Boolean(errors.email)}
                helperText={errors.email?.message}
                placeholder="name@company.com"
                type="email"
                {...register('email', createEmailValidationRules())}
              />
            </Stack>

            <Box
              sx={{
                columnGap: 1.5,
                display: 'grid',
                gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, minmax(0, 1fr))' },
                rowGap: 2.5
              }}
            >
              <Stack spacing={1}>
                <Stack alignItems="baseline" direction="row" spacing={0.75}>
                  <OverlineText component="span">Organization</OverlineText>
                  <Typography color="text.secondary" component="span" variant="caption">
                    (optional)
                  </Typography>
                </Stack>
                <AppTextField
                  placeholder="Acme Corp"
                  {...register('organization')}
                />
              </Stack>

              <Stack spacing={1}>
                <Stack alignItems="baseline" direction="row" spacing={0.75}>
                  <OverlineText component="span">Occupation</OverlineText>
                  <Typography color="text.secondary" component="span" variant="caption">
                    (optional)
                  </Typography>
                </Stack>
                <AppTextField
                  placeholder="Developer"
                  {...register('occupation')}
                />
              </Stack>
            </Box>

            <Stack spacing={1}>
              <OverlineText>Password</OverlineText>
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
              <OverlineText>Confirm Password</OverlineText>
              <PasswordField
                autoComplete="new-password"
                error={Boolean(errors.confirmPassword)}
                fullWidth
                helperText={errors.confirmPassword?.message}
                placeholder="••••••••"
                {...register('confirmPassword', {
                  required: 'Please confirm your password.',
                  validate: (value, values) =>
                    value === values.password || 'Passwords do not match.'
                })}
              />
            </Stack>

            <PasswordRecommendations password={password} />

            <AppButton
              disabled={!isValid || isGithubLoading}
              fullWidth
              loading={isSubmitting}
              loadingText="Creating Account..."
              endIcon={<ArrowForwardRoundedIcon />}
              type="submit"
              variant="contained"
            >
              Initialize Account
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
