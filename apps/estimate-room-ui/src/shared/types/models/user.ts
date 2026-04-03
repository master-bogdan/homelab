export interface AuthUser {
  readonly avatarUrl: string | null;
  readonly displayName: string;
  readonly email: string;
  readonly id: string;
  readonly occupation: string | null;
  readonly organization: string | null;
}
