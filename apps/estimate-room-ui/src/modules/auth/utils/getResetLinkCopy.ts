import {
  ResetPasswordInvalidLinkCopy,
  ResetPasswordValidationReasons
} from '../constants/resetPasswordValidation';
import type { ResetPasswordValidationReason } from '../types';

export const getResetLinkCopy = (reason?: ResetPasswordValidationReason) =>
  ResetPasswordInvalidLinkCopy[reason ?? ResetPasswordValidationReasons.INVALID];
