import type { ResetPasswordPageStates } from '../constants/resetPasswordPageStates';
import type { ResetPasswordValidationReasons } from '../constants/resetPasswordValidation';

export type ResetPasswordPageState =
  (typeof ResetPasswordPageStates)[keyof typeof ResetPasswordPageStates];

export type ResetPasswordValidationReason =
  (typeof ResetPasswordValidationReasons)[keyof typeof ResetPasswordValidationReasons];

export interface ResetPasswordPayload {
  readonly password: string;
  readonly token: string;
}
