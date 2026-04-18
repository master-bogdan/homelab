import { useForm } from 'react-hook-form';

import { useLoginMutation } from '../store';
import {
  isInvalidCredentialsError,
  resolveApiErrorMessage,
  resolveApiHref
} from '../utils';

import {
  useAuthContinuation,
  useFormRootError,
  useGithubAuthRedirect
} from './index';
import type { LoginFormValues } from '../types';

export const useLoginPage = () => {
  const { createPendingRequest } = useAuthContinuation();
  const [login] = useLoginMutation();
  const form = useForm<LoginFormValues>({
    mode: 'onChange',
    defaultValues: {
      email: '',
      password: ''
    },
    reValidateMode: 'onChange'
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

  const submit = form.handleSubmit(async (values) => {
    form.clearErrors();

    let pendingRequest;

    try {
      pendingRequest = await createPendingRequest();
    } catch (error) {
      const message = resolveApiErrorMessage(error, 'Unable to sign in right now.');

      setRootError(message);
      return;
    }

    const result = await login({
      continue: pendingRequest.continueUrl,
      email: values.email,
      password: values.password
    });

    if (result.error) {
      const message = isInvalidCredentialsError(result.error)
        ? 'Email or password is incorrect.'
        : resolveApiErrorMessage(result.error, 'Unable to sign in right now.');

      setRootError(message);
      return;
    }

    window.location.assign(resolveApiHref(pendingRequest.continueUrl));
  });

  return {
    form,
    isGithubLoading,
    onSubmit: submit,
    onSubmitWithGithub: startGithubRedirect
  };
};
