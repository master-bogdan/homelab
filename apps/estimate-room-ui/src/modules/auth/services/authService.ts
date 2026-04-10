import type { AppDispatch } from '@/app/store/store';
import { accessTokenStorage } from '@/shared/api';

import { createApiUrl } from '../utils';
import type {
  ForgotPasswordPayload,
  LoginPayload,
  RegisterPayload,
  ResetPasswordPayload
} from '../types';

import { authApi } from './authApi';

export const authService = {
  exchangeAuthorizationCode: async (
    dispatch: AppDispatch,
    {
      clientId,
      code,
      codeVerifier,
      redirectUri
    }: {
      readonly clientId: string;
      readonly code: string;
      readonly codeVerifier: string;
      readonly redirectUri: string;
    }
  ) =>
    dispatch(
      authApi.endpoints.exchangeAuthorizationCode.initiate({
        clientId,
        code,
        codeVerifier,
        redirectUri
      })
    ).unwrap(),
  fetchSession: async (dispatch: AppDispatch) =>
    dispatch(authApi.endpoints.fetchSession.initiate(undefined, {
      forceRefetch: true,
      subscribe: false
    })).unwrap(),
  forgotPassword: async (dispatch: AppDispatch, payload: ForgotPasswordPayload) =>
    dispatch(authApi.endpoints.forgotPassword.initiate(payload)).unwrap(),
  getGithubLoginUrl: (continueUrl: string) =>
    createApiUrl('auth/github/login', { continue: continueUrl }).toString(),
  hasStoredAccessToken: () => Boolean(accessTokenStorage.get()),
  login: async (dispatch: AppDispatch, payload: LoginPayload) =>
    dispatch(authApi.endpoints.login.initiate(payload)).unwrap(),
  logout: async (dispatch: AppDispatch) => {
    try {
      return await dispatch(authApi.endpoints.logout.initiate()).unwrap();
    } finally {
      accessTokenStorage.clear();
    }
  },
  refreshAccessToken: async (dispatch: AppDispatch) =>
    dispatch(authApi.endpoints.refreshAccessToken.initiate()).unwrap(),
  register: async (dispatch: AppDispatch, payload: RegisterPayload) =>
    dispatch(authApi.endpoints.register.initiate(payload)).unwrap(),
  resetPassword: async (dispatch: AppDispatch, payload: ResetPasswordPayload) =>
    dispatch(authApi.endpoints.resetPassword.initiate(payload)).unwrap(),
  validateResetPasswordToken: async (dispatch: AppDispatch, token: string) =>
    dispatch(authApi.endpoints.validateResetPasswordToken.initiate(token, {
      forceRefetch: true,
      subscribe: false
    })).unwrap()
};
