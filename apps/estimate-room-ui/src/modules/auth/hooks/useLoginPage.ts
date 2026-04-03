import { usePageTitle } from '@/shared/hooks';

export const useLoginPage = () => {
  usePageTitle('Login');

  return {
    readinessItems: [
      'Redux auth state is wired and ready to hydrate from the backend session.',
      'Shared API and WebSocket clients are isolated from components.',
      'Protected routing is already redirecting unauthenticated sessions here.'
    ]
  };
};
