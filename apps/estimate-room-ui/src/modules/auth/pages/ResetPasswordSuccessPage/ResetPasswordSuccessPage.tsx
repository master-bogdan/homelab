import { AuthNarrowCard, AuthPageLayout } from '../../components';
import { ResetPasswordSuccessContent } from './components';

export const ResetPasswordSuccessPage = () => (
  <AuthPageLayout>
    <AuthNarrowCard>
      <ResetPasswordSuccessContent />
    </AuthNarrowCard>
  </AuthPageLayout>
);
