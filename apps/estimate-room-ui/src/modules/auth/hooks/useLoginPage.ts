import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { useAppDispatch } from '@/shared/store';

import { authService } from '../services';
import { setSession } from '../store';
import { isInvalidCredentialsError, resolveApiErrorMessage, resolveApiHref } from '../utils';

import { useAuthContinuation } from './useAuthContinuation';

interface LoginFormValues {
  readonly email: string;
  readonly password: string;
}

export const useLoginPage = () => {
  const dispatch = useAppDispatch();
  const { createPendingRequest } = useAuthContinuation();
  const [isGithubLoading, setIsGithubLoading] = useState(false);
  const form = useForm<LoginFormValues>({
    mode: 'onChange',
    defaultValues: {
      email: '',
      password: ''
    },
    reValidateMode: 'onChange'
  });

  const submit = form.handleSubmit(async (values) => {
    form.clearErrors();

    try {
      const pendingRequest = await createPendingRequest();
      const user = await authService.login(dispatch, {
        continue: pendingRequest.continueUrl,
        email: values.email,
        password: values.password
      });

      dispatch(setSession(user));
      window.location.assign(resolveApiHref(pendingRequest.continueUrl));
    } catch (error) {
      const message = isInvalidCredentialsError(error)
        ? 'Email or password is incorrect.'
        : resolveApiErrorMessage(error, 'Unable to sign in right now.');

      form.setError('root', {
        message,
        type: 'server'
      });
    }
  });

  const signInWithGithub = async () => {
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
    onSubmitWithGithub: signInWithGithub
  };
};
