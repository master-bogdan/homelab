import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { useAppDispatch } from '@/shared/store';

import { submitLogin } from '../store';
import { createGithubLoginUrl, resolveApiErrorMessage, resolveApiHref } from '../utils';

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
      await dispatch(submitLogin({
        continue: pendingRequest.continueUrl,
        email: values.email,
        password: values.password
      })).unwrap();

      window.location.assign(resolveApiHref(pendingRequest.continueUrl));
    } catch (error) {
      const message = resolveApiErrorMessage(error, 'Unable to sign in right now.');

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

      window.location.assign(createGithubLoginUrl(pendingRequest.continueUrl));
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
