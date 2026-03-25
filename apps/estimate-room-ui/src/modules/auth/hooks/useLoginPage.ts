import { useForm } from 'react-hook-form';

import { useAppDispatch } from '@/app/store/hooks';

import { authService } from '../services';
import { setSession } from '../store';
import { isInvalidCredentialsError, resolveApiErrorMessage } from '../utils';

import { useAuthContinuation } from './useAuthContinuation';

interface LoginFormValues {
  readonly email: string;
  readonly password: string;
}

export const useLoginPage = () => {
  const dispatch = useAppDispatch();
  const { createPendingRequest } = useAuthContinuation();
  const form = useForm<LoginFormValues>({
    mode: 'onBlur',
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
      const user = await authService.login({
        continue: pendingRequest.continueUrl,
        email: values.email,
        password: values.password
      });

      dispatch(setSession(user));
      window.location.assign(pendingRequest.continueUrl);
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
    form.clearErrors();

    try {
      const pendingRequest = await createPendingRequest();

      window.location.assign(authService.getGithubLoginUrl(pendingRequest.continueUrl));
    } catch (error) {
      form.setError('root', {
        message: resolveApiErrorMessage(error, 'Unable to start GitHub sign-in.'),
        type: 'server'
      });
    }
  };

  return {
    form,
    onSubmit: submit,
    onSubmitWithGithub: signInWithGithub
  };
};
