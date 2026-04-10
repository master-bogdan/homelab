import { useEffect, useState } from 'react';
import { useForm, useWatch } from 'react-hook-form';

import { useAppDispatch } from '@/app/store/hooks';

import { authService } from '../services';
import { setSession } from '../store';
import { isEmailAlreadyInUseError, resolveApiErrorMessage, resolveApiHref } from '../utils';

import { useAuthContinuation } from './useAuthContinuation';

interface RegisterFormValues {
  readonly confirmPassword: string;
  readonly displayName: string;
  readonly email: string;
  readonly occupation: string;
  readonly organization: string;
  readonly password: string;
}

const normalizeOptionalField = (value: string) => {
  const trimmedValue = value.trim();

  return trimmedValue ? trimmedValue : undefined;
};

export const useRegisterPage = () => {
  const dispatch = useAppDispatch();
  const { createPendingRequest } = useAuthContinuation();
  const [isGithubLoading, setIsGithubLoading] = useState(false);
  const form = useForm<RegisterFormValues>({
    mode: 'onChange',
    defaultValues: {
      confirmPassword: '',
      displayName: '',
      email: '',
      occupation: '',
      organization: '',
      password: ''
    },
    reValidateMode: 'onChange'
  });
  const password = useWatch({
    control: form.control,
    name: 'password'
  });
  const confirmPasswordTouched = form.formState.touchedFields.confirmPassword;

  useEffect(() => {
    if (!confirmPasswordTouched) {
      return;
    }

    void form.trigger('confirmPassword');
  }, [confirmPasswordTouched, form, password]);

  const submit = form.handleSubmit(async ({ confirmPassword: _confirmPassword, ...values }) => {
    form.clearErrors();

    try {
      const pendingRequest = await createPendingRequest();
      const user = await authService.register(dispatch, {
        continue: pendingRequest.continueUrl,
        displayName: values.displayName.trim(),
        email: values.email,
        occupation: normalizeOptionalField(values.occupation),
        organization: normalizeOptionalField(values.organization),
        password: values.password
      });

      dispatch(setSession(user));
      window.location.assign(resolveApiHref(pendingRequest.continueUrl));
    } catch (error) {
      if (isEmailAlreadyInUseError(error)) {
        form.setError('email', {
          message: 'This email is already registered.',
          type: 'server'
        });
        return;
      }

      form.setError('root', {
        message: resolveApiErrorMessage(error, 'Unable to create your account right now.'),
        type: 'server'
      });
    }
  });

  const signUpWithGithub = async () => {
    if (isGithubLoading) {
      return;
    }

    form.clearErrors();
    setIsGithubLoading(true);

    try {
      const pendingRequest = await createPendingRequest();

      window.location.assign(authService.getGithubLoginUrl(pendingRequest.continueUrl));
    } catch (error) {
      setIsGithubLoading(false);
      form.setError('root', {
        message: resolveApiErrorMessage(error, 'Unable to start GitHub sign-in.'),
        type: 'server'
      });
    }
  };

  return {
    form,
    isGithubLoading,
    onSubmit: submit,
    onSubmitWithGithub: signUpWithGithub,
    password
  };
};
