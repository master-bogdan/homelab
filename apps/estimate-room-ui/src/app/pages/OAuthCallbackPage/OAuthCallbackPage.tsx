import { AuthNarrowCard, AuthPageLayout, OAuthCallbackState } from '@/modules/auth/components';
import { useOAuthCallbackPage } from '@/modules/auth/hooks';

export const OAuthCallbackPage = () => {
  const { errorMessage, isLoading } = useOAuthCallbackPage();

  return (
    <AuthPageLayout>
      <AuthNarrowCard>
        <OAuthCallbackState errorMessage={errorMessage} isLoading={isLoading} />
      </AuthNarrowCard>
    </AuthPageLayout>
  );
};
