import type { FormEventHandler } from 'react';
import type { UseFormReturn } from 'react-hook-form';

import {
  AppAlert,
  AppBox,
  AppButton,
  AppStack,
  OverlineText
} from '@/shared/ui';

import { PasswordField, PasswordRecommendations } from '../../../../components';
import { createPasswordValidationRules } from '../../../../utils';
import type { ResetPasswordFormValues } from '../../types';

interface ResetPasswordFormProps {
  readonly form: UseFormReturn<ResetPasswordFormValues>;
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
  readonly password: string;
}

export const ResetPasswordForm = ({
  form,
  onSubmit,
  password
}: ResetPasswordFormProps) => {
  const {
    formState: { errors, isSubmitting, isValid },
    register
  } = form;

  return (
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
  );
};
