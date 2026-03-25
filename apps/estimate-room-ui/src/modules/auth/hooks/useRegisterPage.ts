import { useForm } from 'react-hook-form';

import { useAppDispatch } from '@/app/store/hooks';

import { authService } from '../services';
import { setSession } from '../store';
import { isEmailAlreadyInUseError, resolveApiErrorMessage } from '../utils';

import { useAuthContinuation } from './useAuthContinuation';

interface RegisterFormValues {
  readonly displayName: string;
  readonly email: string;
  readonly password: string;
}

export const useRegisterPage = () => {
  const dispatch = useAppDispatch();
  const { createPendingRequest } = useAuthContinuation();
  const form = useForm<RegisterFormValues>({
    defaultValues: {
      displayName: '',
      email: '',
      password: ''
    }
  });

  const submit = form.handleSubmit(async (values) => {
    form.clearErrors();

    try {
      const pendingRequest = await createPendingRequest();
      const user = await authService.register({
        continue: pendingRequest.continueUrl,
        displayName: values.displayName.trim(),
        email: values.email,
        password: values.password
      });

      dispatch(setSession(user));
      window.location.assign(pendingRequest.continueUrl);
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
    onSubmitWithGithub: signUpWithGithub
  };
};
