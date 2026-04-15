import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { useForgotPasswordMutation } from '../store';
import { resolveApiErrorMessage } from '../utils';

interface ForgotPasswordFormValues {
  readonly email: string;
}

export const useForgotPasswordPage = () => {
  const [forgotPassword] = useForgotPasswordMutation();
  const [submittedEmail, setSubmittedEmail] = useState<string | null>(null);
  const [isResending, setIsResending] = useState(false);
  const form = useForm<ForgotPasswordFormValues>({
    mode: 'onChange',
    defaultValues: {
      email: ''
    },
    reValidateMode: 'onChange'
  });

  const submit = form.handleSubmit(async ({ email }) => {
    form.clearErrors();

    try {
      await forgotPassword({ email }).unwrap();
      setSubmittedEmail(email);
    } catch (error) {
      form.setError('root', {
        message: resolveApiErrorMessage(error, 'Unable to send a reset link right now.'),
        type: 'server'
      });
    }
  });

  const resend = async () => {
    if (!submittedEmail || isResending) {
      return;
    }

    setIsResending(true);

    try {
      await forgotPassword({ email: submittedEmail }).unwrap();
    } finally {
      setIsResending(false);
    }
  };

  return {
    form,
    isResending,
    isSubmitted: submittedEmail !== null,
    onResend: resend,
    onSubmit: submit,
    submittedEmail
  };
};
