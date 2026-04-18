import { AppRoutes } from '@/app/router/routePaths';

import { AuthPageFooter } from '../AuthPageFooter';

export const RegisterFooter = () => (
  <AuthPageFooter
    linkLabel="Sign In"
    prompt="Already have an account?"
    to={AppRoutes.LOGIN}
  />
);
