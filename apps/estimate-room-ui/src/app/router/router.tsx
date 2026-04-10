import { Navigate, createBrowserRouter } from 'react-router-dom';

import { AuthLayout, DashboardLayout } from '@/app/layouts';
import {
  ForgotPasswordPage,
  LoginPage,
  OAuthCallbackPage,
  RegisterPage,
  ResetPasswordPage,
  ResetPasswordSuccessPage
} from '@/modules/auth';
import { DashboardPage } from '@/modules/dashboard';
import { HistoryPage, HistoryRoomPage } from '@/modules/history';
import { ProfilePage } from '@/modules/profile';
import { NewRoomPage, RoomDetailsPage } from '@/modules/rooms';
import { SettingsPage } from '@/modules/settings';
import { TeamDetailsPage } from '@/modules/teams';
import { appRoutes } from '@/shared/constants/routes';

import { NotFoundPage } from './NotFoundPage';
import { ProtectedRoute } from './ProtectedRoute';

export const router = createBrowserRouter([
  {
    path: appRoutes.root,
    element: <Navigate replace to={appRoutes.dashboard} />
  },
  {
    element: <AuthLayout />,
    children: [
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
    ]
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <DashboardLayout />,
        children: [
          {
            path: appRoutes.dashboard,
            element: <DashboardPage />
          },
          {
            path: appRoutes.roomsNew,
            element: <NewRoomPage />
          },
          {
            path: appRoutes.roomDetails,
            element: <RoomDetailsPage />
          },
          {
            path: appRoutes.history,
            element: <HistoryPage />
          },
          {
            path: appRoutes.historyRoom,
            element: <HistoryRoomPage />
          },
          {
            path: appRoutes.teamDetails,
            element: <TeamDetailsPage />
          },
          {
            path: appRoutes.profile,
            element: <ProfilePage />
          },
          {
            path: appRoutes.settings,
            element: <SettingsPage />
          }
        ]
      }
    ]
  },
  {
    path: '*',
    element: <NotFoundPage />
  }
]);
