import { useState } from 'react';

import type { PendingAuthorizationRequest } from '../types';
import {
  createGithubLoginUrl,
  resolveApiErrorMessage
} from '../utils';

interface UseGithubAuthRedirectOptions {
  readonly clearErrors: () => void;
  readonly createPendingRequest: () => Promise<PendingAuthorizationRequest>;
  readonly fallbackMessage: string;
  readonly setRootError: (message: string) => void;
}

export const useGithubAuthRedirect = ({
  clearErrors,
  createPendingRequest,
  fallbackMessage,
  setRootError
}: UseGithubAuthRedirectOptions) => {
  const [isGithubLoading, setIsGithubLoading] = useState(false);

  const startGithubRedirect = () => {
    if (isGithubLoading) {
      return;
    }

    clearErrors();
    setIsGithubLoading(true);

    createPendingRequest()
      .then((pendingRequest) => {
        window.location.assign(createGithubLoginUrl(pendingRequest.continueUrl));
      })
      .catch((error) => {
        setIsGithubLoading(false);
        setRootError(resolveApiErrorMessage(error, fallbackMessage));
      });
  };

  return {
    isGithubLoading,
    startGithubRedirect
  };
};
