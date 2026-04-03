import { apiClient } from '@/shared/api';
import type { AuthUser } from '@/shared/types';

import type { LoginPayload } from '../types';

export const authService = {
  fetchSession: async () => Promise.resolve<AuthUser | null>(null),
  login: async (payload: LoginPayload) => apiClient.post<AuthUser>('auth/login', payload),
  logout: async () => Promise.resolve()
};
