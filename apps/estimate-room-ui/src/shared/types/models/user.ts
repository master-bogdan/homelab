export interface AuthUser {
  readonly displayName: string;
  readonly email: string;
  readonly id: string;
  readonly role: 'admin' | 'member';
  readonly teamIds: string[];
}
