import ArrowForwardRoundedIcon from '@mui/icons-material/ArrowForwardRounded';
import type { FormEventHandler } from 'react';
import type { UseFormReturn } from 'react-hook-form';

import {
  AppAlert,
  AppBox,
  AppButton,
  AppStack,
  AppTextField,
  AppTypography,
  OverlineText
} from '@/shared/components';

import { createEmailValidationRules, createPasswordValidationRules } from '../../utils';
import { registerPageOptionalFieldsSx } from './styles';
import type { RegisterFormValues } from '../../types';
import { AuthActionDivider } from '../AuthActionDivider';
import { AuthGithubButton } from '../AuthGithubButton';
import { PasswordField } from '../PasswordField';
import { PasswordRecommendations } from '../PasswordRecommendations';

interface RegisterFormProps {
  readonly form: UseFormReturn<RegisterFormValues>;
  readonly isGithubLoading: boolean;
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
  readonly onSubmitWithGithub: () => void;
  readonly password: string;
}

export const RegisterForm = ({
  form,
  isGithubLoading,
  onSubmit,
  onSubmitWithGithub,
  password
}: RegisterFormProps) => {
  const {
    formState: { errors, isSubmitting, isValid },
    register
  } = form;

  return (
    <AppBox component="form" noValidate onSubmit={onSubmit}>
      <AppStack spacing={2.5}>
        {errors.root?.message ? <AppAlert severity="error">{errors.root.message}</AppAlert> : null}

        <AppStack spacing={1}>
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
        </AppStack>

        <AppStack spacing={1}>
          <OverlineText>Work Email</OverlineText>
          <AppTextField
            autoComplete="email"
            error={Boolean(errors.email)}
            helperText={errors.email?.message}
            placeholder="name@company.com"
            type="email"
            {...register('email', createEmailValidationRules())}
          />
        </AppStack>

        <AppBox sx={registerPageOptionalFieldsSx}>
          <AppStack spacing={1}>
            <AppStack alignItems="baseline" direction="row" spacing={0.75}>
              <OverlineText component="span">Organization</OverlineText>
              <AppTypography color="text.secondary" component="span" variant="caption">
                (optional)
              </AppTypography>
            </AppStack>
            <AppTextField placeholder="Acme Corp" {...register('organization')} />
          </AppStack>

          <AppStack spacing={1}>
            <AppStack alignItems="baseline" direction="row" spacing={0.75}>
              <OverlineText component="span">Occupation</OverlineText>
              <AppTypography color="text.secondary" component="span" variant="caption">
                (optional)
              </AppTypography>
            </AppStack>
            <AppTextField placeholder="Developer" {...register('occupation')} />
          </AppStack>
        </AppBox>

        <AppStack spacing={1}>
          <OverlineText>Password</OverlineText>
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
        </AppStack>

        <PasswordRecommendations password={password} />

        <AppButton
          disabled={!isValid || isGithubLoading}
          endIcon={<ArrowForwardRoundedIcon />}
          fullWidth
          loading={isSubmitting}
          loadingText="Creating Account..."
          type="submit"
          variant="contained"
        >
          Initialize Account
        </AppButton>

        <AuthActionDivider />

        <AuthGithubButton
          disabled={isSubmitting}
          loading={isGithubLoading}
          onClick={onSubmitWithGithub}
        >
          Sign up with GitHub
        </AuthGithubButton>
      </AppStack>
    </AppBox>
  );
};
