import { useMemo } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { AppRoutes } from '@/app/router/routePaths';

import {
  ResetPasswordErrorMessages,
  ResetPasswordSearchParams
} from '../constants/resetPasswordValidation';
import { useResetPasswordMutation } from '../store';
import { getResetLinkCopy } from '../utils/getResetLinkCopy';
import { resolveApiErrorMessage } from '../utils/errorMessages';
import { resolveResetPasswordValidationReason } from '../utils/resetPasswordValidation';
import { useFormRootError } from './useFormRootError';
import { useResetPasswordForm } from './useResetPasswordForm';
import { useResetPasswordTokenValidation } from './useResetPasswordTokenValidation';

export const useResetPasswordPage = () => {
  const navigate = useNavigate();
  const [resetPassword] = useResetPasswordMutation();
  const [searchParams] = useSearchParams();
  const token = searchParams.get(ResetPasswordSearchParams.TOKEN)?.trim() ?? '';
  const { form, password } = useResetPasswordForm();
  const { setRootError } = useFormRootError(form);
  const {
    markResetLinkInvalid,
    pageError,
    pageState,
    validationReason
  } = useResetPasswordTokenValidation(token);
  const invalidLinkCopy = useMemo(
    () => getResetLinkCopy(validationReason),
    [validationReason]
  );

  const submit = form.handleSubmit(async ({ password: nextPassword }) => {
    form.clearErrors();

    if (!token) {
      markResetLinkInvalid();
      return;
    }

    const result = await resetPassword({
      password: nextPassword,
      token
    });

    if (result.error) {
      const message = resolveApiErrorMessage(
        result.error,
        ResetPasswordErrorMessages.RESET_FAILED
      );

      const resetPasswordValidationReason =
        resolveResetPasswordValidationReason(message);

      if (resetPasswordValidationReason) {
        markResetLinkInvalid(resetPasswordValidationReason);
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
