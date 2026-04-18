import { AppPageState } from '@/shared/components';

export const ResetPasswordValidatingState = () => (
  <AppPageState
    description="Confirming your password reset link before you choose a new password."
    isLoading
    title="Validating Link"
  />
);
