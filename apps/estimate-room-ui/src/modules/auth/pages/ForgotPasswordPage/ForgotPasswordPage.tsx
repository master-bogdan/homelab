import { AuthCard, AuthPageHeader, AuthPageLayout } from '../../components';
import {
  ForgotPasswordForm,
  ForgotPasswordResendPrompt,
  ForgotPasswordSubmittedContent
} from './components';
import { useForgotPasswordPage } from './hooks';

export const ForgotPasswordPage = () => {
  const {
    errorMessage,
    form,
    isResending,
    isSubmitted,
    onSubmit,
    resend
  } = useForgotPasswordPage();

  return (
    <AuthPageLayout>
      <AuthPageHeader
        description={
          isSubmitted
            ? "We've sent a password reset link to your email address."
            : 'Enter your work email and we will send you a secure reset link.'
        }
        title={isSubmitted ? 'Check Your Email' : 'Reset Password'}
      />
      <AuthCard>
        {isSubmitted ? (
          <ForgotPasswordSubmittedContent />
        ) : (
          <ForgotPasswordForm
            errorMessage={errorMessage}
            form={form}
            onSubmit={onSubmit}
          />
        )}
      </AuthCard>
      {isSubmitted ? (
        <ForgotPasswordResendPrompt isResending={isResending} onResend={resend} />
      ) : null}
    </AuthPageLayout>
  );
};
