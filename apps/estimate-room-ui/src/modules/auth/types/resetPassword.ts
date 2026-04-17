export type ResetPasswordValidationReason = 'expired' | 'invalid' | 'used';

export interface ResetPasswordPayload {
  readonly password: string;
  readonly token: string;
}
