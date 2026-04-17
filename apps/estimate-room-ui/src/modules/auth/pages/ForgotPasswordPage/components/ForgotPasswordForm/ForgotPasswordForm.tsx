import type { FormEventHandler } from 'react';
import type { UseFormReturn } from 'react-hook-form';
import {
  AppAlert,
  AppBox,
  AppButton,
  AppStack,
  AppTextField,
  OverlineText
} from '@/shared/ui';

import { AuthBackToSignInLink } from '../../../../components';
import { createEmailValidationRules } from '../../../../utils';
import type { ForgotPasswordFormValues } from '../../types';

interface ForgotPasswordFormProps {
  readonly errorMessage: string | null;
  readonly form: UseFormReturn<ForgotPasswordFormValues>;
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
}

export const ForgotPasswordForm = ({
  errorMessage,
  form,
  onSubmit
}: ForgotPasswordFormProps) => {
  const {
    formState: { errors, isSubmitting, isValid },
    register
  } = form;

  return (
    <AppBox component="form" noValidate onSubmit={onSubmit}>
      <AppStack spacing={2.5}>
        {errorMessage ? <AppAlert severity="error">{errorMessage}</AppAlert> : null}
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
        <AuthBackToSignInLink placement="form" />
      </AppStack>
    </AppBox>
  );
};
