import { ResetPasswordValidationReasons } from '../../constants/resetPasswordValidation';
import { getResetLinkCopy } from '../getResetLinkCopy';

describe('getResetLinkCopy', () => {
  it('returns reset link copy for each validation reason', () => {
    expect(getResetLinkCopy(ResetPasswordValidationReasons.EXPIRED).title).toBe(
      'Expired Link'
    );
    expect(getResetLinkCopy(ResetPasswordValidationReasons.USED).title).toBe(
      'Link Already Used'
    );
    expect(getResetLinkCopy(ResetPasswordValidationReasons.INVALID).title).toBe(
      'Invalid Link'
    );
  });

  it('uses invalid link copy as the default', () => {
    expect(getResetLinkCopy()).toEqual(
      getResetLinkCopy(ResetPasswordValidationReasons.INVALID)
    );
  });
});
