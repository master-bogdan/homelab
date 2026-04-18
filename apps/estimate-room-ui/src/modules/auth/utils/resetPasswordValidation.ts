import {
  ResetPasswordValidationMessageFragments,
  ResetPasswordValidationReasons
} from '../constants/resetPasswordValidation';
import type { ResetPasswordValidationReason } from '../types/resetPassword';

const resetPasswordValidationReasonEntries = [
  {
    fragment: ResetPasswordValidationMessageFragments.EXPIRED,
    reason: ResetPasswordValidationReasons.EXPIRED
  },
  {
    fragment: ResetPasswordValidationMessageFragments.USED,
    reason: ResetPasswordValidationReasons.USED
  },
  {
    fragment: ResetPasswordValidationMessageFragments.INVALID,
    reason: ResetPasswordValidationReasons.INVALID
  }
] as const;

export const resolveResetPasswordValidationReason = (
  message: string
): ResetPasswordValidationReason | null => {
  const normalizedMessage = message.toLowerCase();
  const matchedReason = resetPasswordValidationReasonEntries.find(({ fragment }) =>
    normalizedMessage.includes(fragment)
  );

  return matchedReason?.reason ?? null;
};
