import {
  AuthCard,
  AuthPageHeader,
  AuthPageLayout,
  LoginFooter,
  LoginForm
} from '@/modules/auth/components';
import { useLoginPage } from '@/modules/auth/hooks';

export const LoginPage = () => {
  const {
    form,
    isGithubLoading,
    onSubmit,
    onSubmitWithGithub
  } = useLoginPage();

  return (
    <AuthPageLayout>
      <AuthPageHeader
        description="Access your precision estimation workspace."
        title="Welcome Back"
      />
      <AuthCard>
        <LoginForm
          form={form}
          isGithubLoading={isGithubLoading}
          onSubmit={onSubmit}
          onSubmitWithGithub={onSubmitWithGithub}
        />
      </AuthCard>
      <LoginFooter />
    </AuthPageLayout>
  );
};
