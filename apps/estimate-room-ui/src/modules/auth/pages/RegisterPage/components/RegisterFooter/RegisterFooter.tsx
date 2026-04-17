import { AppRoutes } from '@/shared/constants/routes';

import { AuthPageFooter } from '../../../../components';

export const RegisterFooter = () => (
  <AuthPageFooter
    linkLabel="Sign In"
    prompt="Already have an account?"
    to={AppRoutes.LOGIN}
  />
);
