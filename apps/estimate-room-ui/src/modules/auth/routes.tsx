import type { RouteObject } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';

import {
  ForgotPasswordPage,
  LoginPage,
  OAuthCallbackPage,
  RegisterPage,
  ResetPasswordPage,
  ResetPasswordSuccessPage
} from './pages';

export const AuthRoutes: RouteObject[] = [
  {
    path: AppRoutes.LOGIN,
    element: <LoginPage />
  },
  {
    path: AppRoutes.REGISTER,
    element: <RegisterPage />
  },
  {
    path: AppRoutes.FORGOT_PASSWORD,
    element: <ForgotPasswordPage />
  },
  {
    path: AppRoutes.RESET_PASSWORD,
    element: <ResetPasswordPage />
  },
  {
    path: AppRoutes.RESET_PASSWORD_SUCCESS,
    element: <ResetPasswordSuccessPage />
  },
  {
    path: AppRoutes.AUTH_CALLBACK,
    element: <OAuthCallbackPage />
  }
];
