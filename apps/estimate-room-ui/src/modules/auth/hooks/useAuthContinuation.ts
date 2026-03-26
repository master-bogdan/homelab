import { useLocation } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';

import { ensurePendingAuthorizationRequest } from '../utils';

interface RedirectStateLike {
  readonly from?: {
    readonly hash: string;
    readonly pathname: string;
    readonly search: string;
  };
}

const resolveRedirectTarget = (state: RedirectStateLike | null) => {
  const from = state?.from;

  if (!from) {
    return appRoutes.dashboard;
  }

  return `${from.pathname}${from.search}${from.hash}`;
};

export const useAuthContinuation = () => {
  const location = useLocation();
  const continueUrl = new URLSearchParams(location.search).get('continue');
  const redirectTo = resolveRedirectTarget(location.state as RedirectStateLike | null);

  const createPendingRequest = async () =>
    ensurePendingAuthorizationRequest(redirectTo, continueUrl);

  return {
    continueUrl,
    createPendingRequest,
    redirectTo
  };
};
