import { AppRoutes } from '@/shared/constants/routes';

import { AuthPageFooter } from '../../../../components';

export const LoginFooter = () => (
  <AuthPageFooter
    linkLabel="Register now"
    prompt="Don't have an account?"
    to={AppRoutes.REGISTER}
  />
);
