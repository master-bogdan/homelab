import { AppRoutes } from '@/app/router/routePaths';

import { AuthPageFooter } from '../AuthPageFooter';

export const ResetPasswordFooter = () => (
  <AuthPageFooter
    linkLabel="Back to Login"
    prompt="Remember your password?"
    to={AppRoutes.LOGIN}
  />
);
