import {
  AuthCard,
  AuthPageHeader,
  AuthPageLayout,
  ResetPasswordFooter,
  ResetPasswordForm,
  ResetPasswordInvalidState,
  ResetPasswordValidatingState
} from '@/modules/auth/components';
import { ResetPasswordPageStates } from '@/modules/auth';
import { useResetPasswordPage } from '@/modules/auth/hooks';

export const ResetPasswordPage = () => {
  const {
    form,
    invalidLinkCopy,
    onSubmit,
    pageError,
    pageState,
    password
  } = useResetPasswordPage();

  if (pageState === ResetPasswordPageStates.INVALID) {
    return (
      <ResetPasswordInvalidState
        description={invalidLinkCopy.description}
        pageError={pageError}
        title={invalidLinkCopy.title}
      />
    );
  }

  return (
    <AuthPageLayout>
      <AuthPageHeader
        description="Create a strong new password for your EstimateRoom account."
        title={
          pageState === ResetPasswordPageStates.VALIDATING
            ? 'Validating Link'
            : 'Set New Password'
        }
      />
      <AuthCard>
        {pageState === ResetPasswordPageStates.VALIDATING ? (
          <ResetPasswordValidatingState />
        ) : (
          <ResetPasswordForm form={form} onSubmit={onSubmit} password={password} />
        )}
      </AuthCard>
      <ResetPasswordFooter />
    </AuthPageLayout>
  );
};
