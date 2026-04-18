import { useForm, useWatch } from 'react-hook-form';

import type { ResetPasswordFormValues } from '../types';
import { useConfirmPasswordRevalidation } from './useConfirmPasswordRevalidation';

export const useResetPasswordForm = () => {
  const form = useForm<ResetPasswordFormValues>({
    mode: 'onChange',
    defaultValues: {
      confirmPassword: '',
      password: ''
    },
    reValidateMode: 'onChange'
  });
  const password = useWatch({
    control: form.control,
    name: 'password'
  });

  useConfirmPasswordRevalidation({
    confirmPasswordField: 'confirmPassword',
    form,
    password
  });

  return {
    form,
    password
  };
};
