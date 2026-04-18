import {
  AuthNarrowCard,
  AuthPageLayout,
  ResetPasswordSuccessContent
} from '@/modules/auth/components';

export const ResetPasswordSuccessPage = () => (
  <AuthPageLayout>
    <AuthNarrowCard>
      <ResetPasswordSuccessContent />
    </AuthNarrowCard>
  </AuthPageLayout>
);
