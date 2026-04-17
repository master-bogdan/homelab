import { AppRoutes } from '@/shared/constants/routes';

import { AuthPageFooter } from '../../../../components';

export const ResetPasswordFooter = () => (
  <AuthPageFooter
    linkLabel="Back to Login"
    prompt="Remember your password?"
    to={AppRoutes.LOGIN}
  />
);
