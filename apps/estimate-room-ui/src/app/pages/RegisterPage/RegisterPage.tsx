import {
  AuthCard,
  AuthPageHeader,
  AuthPageLayout,
  RegisterFooter,
  RegisterForm
} from '@/modules/auth/components';
import { useRegisterPage } from '@/modules/auth/hooks';

import { registerPageCardSx } from './RegisterPage.styles';

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
