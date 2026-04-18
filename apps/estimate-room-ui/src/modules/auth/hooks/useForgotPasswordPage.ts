import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { useForgotPasswordMutation } from '../store';
import { resolveApiErrorMessage } from '../utils';
import type { ForgotPasswordFormValues } from '../types';

export const useForgotPasswordPage = () => {
  const [forgotPassword, forgotPasswordState] = useForgotPasswordMutation();
  const [submittedEmail, setSubmittedEmail] = useState<string | null>(null);
  const form = useForm<ForgotPasswordFormValues>({
    mode: 'onChange',
    defaultValues: {
      email: ''
    },
    reValidateMode: 'onChange'
  });
  const isSubmitted = submittedEmail !== null;
  const isResending = forgotPasswordState.isLoading;
  const errorMessage = forgotPasswordState.error
    ? resolveApiErrorMessage(
        forgotPasswordState.error,
        'Unable to send a reset link right now.'
      )
    : null;

  const submit = async (email: string) => {
    const result = await forgotPassword({ email });

    if (result.data) {
      setSubmittedEmail(email);
    }
  };

  const onSubmit = form.handleSubmit(async ({ email }) => {
    form.clearErrors();
    await submit(email);
  });

  const resend = async () => {
    if (!submittedEmail || forgotPasswordState.isLoading) {
      return;
    }

    await forgotPassword({ email: submittedEmail });
  };

  return {
    errorMessage,
    form,
    isResending,
    isSubmitted,
    onSubmit,
    resend
  };
};
