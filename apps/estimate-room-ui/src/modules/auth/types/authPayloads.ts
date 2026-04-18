export interface ForgotPasswordPayload {
  readonly email: string;
}

export interface LoginPayload {
  readonly continue: string;
  readonly email: string;
  readonly password: string;
}

export interface RegisterPayload {
  readonly continue: string;
  readonly displayName: string;
  readonly email: string;
  readonly occupation?: string;
  readonly organization?: string;
  readonly password: string;
}
