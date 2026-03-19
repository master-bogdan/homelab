import type { AuthUser } from '@/shared/types';

export interface AuthState {
  readonly status: 'authenticated' | 'unauthenticated' | 'unknown';
  readonly user: AuthUser | null;
}

export interface LoginPayload {
  readonly email: string;
  readonly password: string;
}
