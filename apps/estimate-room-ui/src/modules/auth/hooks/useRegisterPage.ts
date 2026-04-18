import { useForm, useWatch } from 'react-hook-form';

import { useRegisterMutation } from '../store';
import {
  isEmailAlreadyInUseError,
  resolveApiErrorMessage,
  resolveApiHref
} from '../utils';

import {
  useAuthContinuation,
  useConfirmPasswordRevalidation,
  useFormRootError,
  useGithubAuthRedirect
} from './index';
import type { RegisterFormValues } from '../types';

const normalizeOptionalField = (value: string) => {
  const trimmedValue = value.trim();

  return trimmedValue ? trimmedValue : undefined;
};

export const useRegisterPage = () => {
  const { createPendingRequest } = useAuthContinuation();
  const [registerAccount] = useRegisterMutation();
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
  const { setRootError } = useFormRootError(form);
  const {
    isGithubLoading,
    startGithubRedirect
  } = useGithubAuthRedirect({
    clearErrors: form.clearErrors,
    createPendingRequest,
    fallbackMessage: 'Unable to start GitHub sign-in.',
    setRootError
  });

  useConfirmPasswordRevalidation({
    confirmPasswordField: 'confirmPassword',
    form,
    password
  });

  const submit = form.handleSubmit(async (values) => {
    form.clearErrors();

    let pendingRequest;

    try {
      pendingRequest = await createPendingRequest();
    } catch (error) {
      const message = resolveApiErrorMessage(error, 'Unable to create your account right now.');

      if (message === 'This email is already registered.') {
        form.setError('email', {
          message,
          type: 'server'
        });
        return;
      }

      setRootError(message);
      return;
    }

    const result = await registerAccount({
      continue: pendingRequest.continueUrl,
      displayName: values.displayName.trim(),
      email: values.email,
      occupation: normalizeOptionalField(values.occupation),
      organization: normalizeOptionalField(values.organization),
      password: values.password
    });

    if (result.error) {
      const message = isEmailAlreadyInUseError(result.error)
        ? 'This email is already registered.'
        : resolveApiErrorMessage(result.error, 'Unable to create your account right now.');

      if (message === 'This email is already registered.') {
        form.setError('email', {
          message,
          type: 'server'
        });
        return;
      }

      setRootError(message);
      return;
    }

    window.location.assign(resolveApiHref(pendingRequest.continueUrl));
  });

  return {
    form,
    isGithubLoading,
    onSubmit: submit,
    onSubmitWithGithub: startGithubRedirect,
    password
  };
};
