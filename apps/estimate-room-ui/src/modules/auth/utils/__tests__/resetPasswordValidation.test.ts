import { ResetPasswordValidationReasons } from '../../constants/resetPasswordValidation';
import { resolveResetPasswordValidationReason } from '../resetPasswordValidation';

describe('reset password validation utilities', () => {
  it('maps reset token error messages to validation reasons', () => {
    expect(resolveResetPasswordValidationReason('expired reset token')).toBe(
      ResetPasswordValidationReasons.EXPIRED
    );
    expect(resolveResetPasswordValidationReason('used reset token')).toBe(
      ResetPasswordValidationReasons.USED
    );
    expect(resolveResetPasswordValidationReason('invalid reset token')).toBe(
      ResetPasswordValidationReasons.INVALID
    );
  });

  it('returns null when the reset password error is not token-specific', () => {
    expect(resolveResetPasswordValidationReason('backend unavailable')).toBeNull();
  });
});
