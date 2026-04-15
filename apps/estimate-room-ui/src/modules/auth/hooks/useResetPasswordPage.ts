import { useEffect, useMemo, useState } from 'react';
import { useForm, useWatch } from 'react-hook-form';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { useAppDispatch } from '@/shared/store';
import { appRoutes } from '@/shared/constants/routes';

import { authService } from '../services';
import type { ResetPasswordValidationReason } from '../types';
import { clearSession } from '../store';
import { getResetLinkCopy, resolveApiErrorMessage } from '../utils';

interface ResetPasswordFormValues {
  readonly confirmPassword: string;
  readonly password: string;
}

type ResetPasswordPageState = 'invalid' | 'ready' | 'validating';

export const useResetPasswordPage = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
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

  useEffect(() => {
    let isMounted = true;

    if (!token) {
      return;
    }

    const validateToken = async () => {
      try {
        const response = await authService.validateResetPasswordToken(dispatch, token);

        if (!isMounted) {
          return;
        }

        if (response.valid) {
          setPageState('ready');
          return;
        }

        setValidationReason(response.reason ?? 'invalid');
        setPageState('invalid');
      } catch (error) {
        if (!isMounted) {
          return;
        }

        setPageError(resolveApiErrorMessage(error, 'Unable to validate this reset link.'));
        setValidationReason('invalid');
        setPageState('invalid');
      }
    };

    void validateToken();

    return () => {
      isMounted = false;
    };
  }, [token]);

  const password = useWatch({
    control: form.control,
    name: 'password'
  });
  const confirmPasswordTouched = form.formState.touchedFields.confirmPassword;
  const invalidLinkCopy = useMemo(
    () => getResetLinkCopy(validationReason),
    [validationReason]
  );

  useEffect(() => {
    if (!confirmPasswordTouched) {
      return;
    }

    void form.trigger('confirmPassword');
  }, [confirmPasswordTouched, form, password]);

  const submit = form.handleSubmit(async ({ password: nextPassword }) => {
    form.clearErrors();

    if (!token) {
      setPageState('invalid');
      setValidationReason('invalid');
      return;
    }

    try {
      await authService.resetPassword(dispatch, {
        password: nextPassword,
        token
      });

      dispatch(clearSession());
      navigate(appRoutes.resetPasswordSuccess, { replace: true });
    } catch (error) {
      const message = resolveApiErrorMessage(error, 'Unable to reset your password right now.');

      if (
        message.toLowerCase().includes('invalid reset token') ||
        message.toLowerCase().includes('expired reset token') ||
        message.toLowerCase().includes('used reset token')
      ) {
        setValidationReason(
          message.toLowerCase().includes('expired')
            ? 'expired'
            : message.toLowerCase().includes('used')
              ? 'used'
              : 'invalid'
        );
        setPageState('invalid');
        return;
      }

      form.setError('root', {
        message,
        type: 'server'
      });
    }
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
