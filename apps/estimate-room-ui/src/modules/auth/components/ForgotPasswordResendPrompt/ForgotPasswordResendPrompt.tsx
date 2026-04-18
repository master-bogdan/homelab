import { AppLink, AppProgress } from '@/shared/components';

import { forgotPasswordResendLinkSx } from './styles';
import { AuthInlinePrompt } from '../AuthInlinePrompt';

interface ForgotPasswordResendPromptProps {
  readonly isResending: boolean;
  readonly onResend: () => void;
}

export const ForgotPasswordResendPrompt = ({
  isResending,
  onResend
}: ForgotPasswordResendPromptProps) => (
  <AuthInlinePrompt>
    Didn&apos;t receive the email?{' '}
    <AppLink
      color="primary"
      component="button"
      disabled={isResending}
      onClick={onResend}
      sx={forgotPasswordResendLinkSx}
      type="button"
      underline="none"
    >
      {isResending ? <AppProgress color="inherit" size={14} /> : null}
      {isResending ? 'Resending...' : 'Click to resend link'}
    </AppLink>
  </AuthInlinePrompt>
);
