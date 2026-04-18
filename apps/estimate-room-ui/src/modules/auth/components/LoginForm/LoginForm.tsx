import type { FormEventHandler } from 'react';
import type { UseFormReturn } from 'react-hook-form';
import { Link as RouterLink } from 'react-router-dom';

import { AppRoutes } from '@/app/router/routePaths';
import {
  AppAlert,
  AppBox,
  AppButton,
  AppLink,
  AppStack,
  AppTextField,
  OverlineText
} from '@/shared/components';

import { createEmailValidationRules } from '../../utils';
import type { LoginFormValues } from '../../types';
import { AuthActionDivider } from '../AuthActionDivider';
import { AuthGithubButton } from '../AuthGithubButton';
import { PasswordField } from '../PasswordField';

interface LoginFormProps {
  readonly form: UseFormReturn<LoginFormValues>;
  readonly isGithubLoading: boolean;
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
  readonly onSubmitWithGithub: () => void;
}

export const LoginForm = ({
  form,
  isGithubLoading,
  onSubmit,
  onSubmitWithGithub
}: LoginFormProps) => {
  const {
    formState: { errors, isSubmitting, isValid },
    register
  } = form;

  return (
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

        <AuthGithubButton
          disabled={isSubmitting}
          loading={isGithubLoading}
          onClick={onSubmitWithGithub}
        >
          Continue with GitHub
        </AuthGithubButton>
      </AppStack>
    </AppBox>
  );
};
