import { AuthCard, AuthPageHeader, AuthPageLayout } from '../../components';
import { LoginFooter, LoginForm } from './components';
import { useLoginPage } from './hooks';

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
