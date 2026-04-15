import type { RouteObject } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';

import { ForgotPasswordPage } from './ForgotPasswordPage';
import { LoginPage } from './LoginPage';
import { OAuthCallbackPage } from './OAuthCallbackPage';
import { RegisterPage } from './RegisterPage';
import { ResetPasswordPage } from './ResetPasswordPage';
import { ResetPasswordSuccessPage } from './ResetPasswordSuccessPage';

export const authRoutes: RouteObject[] = [
  {
    path: appRoutes.login,
    element: <LoginPage />
  },
  {
    path: appRoutes.register,
    element: <RegisterPage />
  },
  {
    path: appRoutes.forgotPassword,
    element: <ForgotPasswordPage />
  },
  {
    path: appRoutes.resetPassword,
    element: <ResetPasswordPage />
  },
  {
    path: appRoutes.resetPasswordSuccess,
    element: <ResetPasswordSuccessPage />
  },
  {
    path: appRoutes.authCallback,
    element: <OAuthCallbackPage />
  }
];
