import {
  AuthCard,
  AuthPageHeader,
  AuthPageLayout
} from '../../components';
import { RegisterFooter, RegisterForm } from './components';
import { registerPageCardSx } from './styles';
import { useRegisterPage } from './hooks';

export const RegisterPage = () => {
  const {
    form,
    isGithubLoading,
    onSubmit,
    onSubmitWithGithub,
    password
  } = useRegisterPage();

  return (
    <AuthPageLayout>
      <AuthPageHeader
        description="Define your professional profile to start estimating."
        title="Create Workspace"
      />
      <AuthCard sx={registerPageCardSx}>
        <RegisterForm
          form={form}
          isGithubLoading={isGithubLoading}
          onSubmit={onSubmit}
          onSubmitWithGithub={onSubmitWithGithub}
          password={password}
        />
      </AuthCard>
      <RegisterFooter />
    </AuthPageLayout>
  );
};
