import MailOutlineRoundedIcon from '@mui/icons-material/MailOutlineRounded';
import {
  AppBox,
  AppButton,
  AppStack,
  AppTypography
} from '@/shared/components';

import { AuthBackToSignInLink } from '../AuthBackToSignInLink';
import { forgotPasswordSubmittedActionsSx } from './styles';

export const ForgotPasswordSubmittedContent = () => (
  <AppStack spacing={3}>
    <AppTypography align="center" color="text.secondary" variant="body2">
      Please check your inbox and follow the instructions to reset your password.
    </AppTypography>
    <AppButton
      color="secondary"
      component="a"
      fullWidth
      href="mailto:"
      startIcon={<MailOutlineRoundedIcon />}
      variant="contained"
    >
      Open Email App
    </AppButton>
    <AppBox sx={forgotPasswordSubmittedActionsSx}>
      <AppStack alignItems="center" spacing={1.5}>
        <AuthBackToSignInLink />
      </AppStack>
    </AppBox>
  </AppStack>
);
