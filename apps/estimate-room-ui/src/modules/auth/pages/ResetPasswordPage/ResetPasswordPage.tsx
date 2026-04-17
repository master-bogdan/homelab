import {
  AuthCard,
  AuthPageHeader,
  AuthPageLayout
} from '../../components';
import {
  ResetPasswordFooter,
  ResetPasswordForm,
  ResetPasswordInvalidState,
  ResetPasswordValidatingState
} from './components';
import { useResetPasswordPage } from './hooks';

export const ResetPasswordPage = () => {
  const {
    form,
    invalidLinkCopy,
    onSubmit,
    pageError,
    pageState,
    password
  } = useResetPasswordPage();

  if (pageState === 'invalid') {
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
        title={pageState === 'validating' ? 'Validating Link' : 'Set New Password'}
      />
      <AuthCard>
        {pageState === 'validating' ? (
          <ResetPasswordValidatingState />
        ) : (
          <ResetPasswordForm form={form} onSubmit={onSubmit} password={password} />
        )}
      </AuthCard>
      <ResetPasswordFooter />
    </AuthPageLayout>
  );
};
