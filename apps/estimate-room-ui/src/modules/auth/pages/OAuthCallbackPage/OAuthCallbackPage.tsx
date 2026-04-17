import { AuthNarrowCard, AuthPageLayout } from '../../components';
import { OAuthCallbackState } from './components';
import { useOAuthCallbackPage } from './hooks';

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
