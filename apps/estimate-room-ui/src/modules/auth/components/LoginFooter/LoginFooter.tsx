import { AppRoutes } from '@/app/router/routePaths';

import { AuthPageFooter } from '../AuthPageFooter';

export const LoginFooter = () => (
  <AuthPageFooter
    linkLabel="Register now"
    prompt="Don't have an account?"
    to={AppRoutes.REGISTER}
  />
);
