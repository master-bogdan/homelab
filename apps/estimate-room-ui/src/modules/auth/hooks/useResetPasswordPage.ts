import { useEffect, useMemo, useState } from 'react';
import { useForm, useWatch } from 'react-hook-form';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { AppRoutes } from '@/app/router/routePaths';

import type { ResetPasswordValidationReason } from '../types';
import {
  useLazyValidateResetPasswordTokenQuery,
  useResetPasswordMutation
} from '../store';
import { getResetLinkCopy, resolveApiErrorMessage } from '../utils';
import {
  useConfirmPasswordRevalidation,
  useFormRootError
} from './index';
import type { ResetPasswordFormValues } from '../types';

type ResetPasswordPageState = 'invalid' | 'ready' | 'validating';

const resolveResetPasswordValidationReason = (
  message: string
): ResetPasswordValidationReason | null => {
  const normalizedMessage = message.toLowerCase();

  if (normalizedMessage.includes('expired reset token')) {
    return 'expired';
  }

  if (normalizedMessage.includes('used reset token')) {
    return 'used';
  }

  if (normalizedMessage.includes('invalid reset token')) {
    return 'invalid';
  }

  return null;
};

export const useResetPasswordPage = () => {
  const navigate = useNavigate();
  const [validateResetPasswordToken] = useLazyValidateResetPasswordTokenQuery();
  const [resetPassword] = useResetPasswordMutation();
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token')?.trim() ?? '';
  const [pageState, setPageState] = useState<ResetPasswordPageState>(
    token ? 'validating' : 'invalid'
  );
  const [validationReason, setValidationReason] =
    useState<ResetPasswordValidationReason>('invalid');
  const [pageError, setPageError] = useState<string | null>(null);
  const form = useForm<ResetPasswordFormValues>({
    mode: 'onChange',
    defaultValues: {
      confirmPassword: '',
      password: ''
    },
    reValidateMode: 'onChange'
  });
  const { setRootError } = useFormRootError(form);

  useEffect(() => {
    let isMounted = true;

    if (!token) {
      return;
    }

    const validateToken = async () => {
      const response = await validateResetPasswordToken(token, true);

      if (!isMounted) {
        return;
      }

      if (response.data) {
        if (response.data.valid) {
          setPageState('ready');
          return;
        }

        setValidationReason(response.data.reason ?? 'invalid');
        setPageState('invalid');
        return;
      }

      setPageError(resolveApiErrorMessage(response.error, 'Unable to validate this reset link.'));
      setValidationReason('invalid');
      setPageState('invalid');
    };

    validateToken();

    return () => {
      isMounted = false;
    };
  }, [token, validateResetPasswordToken]);

  const password = useWatch({
    control: form.control,
    name: 'password'
  });
  const invalidLinkCopy = useMemo(
    () => getResetLinkCopy(validationReason),
    [validationReason]
  );

  useConfirmPasswordRevalidation({
    confirmPasswordField: 'confirmPassword',
    form,
    password
  });

  const submit = form.handleSubmit(async ({ password: nextPassword }) => {
    form.clearErrors();

    if (!token) {
      setPageState('invalid');
      setValidationReason('invalid');
      return;
    }

    const result = await resetPassword({
      password: nextPassword,
      token
    });

    if (result.error) {
      const message = resolveApiErrorMessage(
        result.error,
        'Unable to reset your password right now.'
      );

      const resetPasswordValidationReason =
        resolveResetPasswordValidationReason(message);

      if (resetPasswordValidationReason) {
        setValidationReason(resetPasswordValidationReason);
        setPageState('invalid');
        return;
      }

      setRootError(message);
      return;
    }

    navigate(AppRoutes.RESET_PASSWORD_SUCCESS, { replace: true });
  });

  return {
    form,
    invalidLinkCopy,
    onSubmit: submit,
    pageError,
    pageState,
    password
  };
};
