import { appRoutes } from '@/shared/constants/routes';

export interface RedirectStateLike {
  readonly from?: {
    readonly hash: string;
    readonly pathname: string;
    readonly search: string;
  };
}

export const resolveAuthRedirectTarget = (state: RedirectStateLike | null) => {
  const from = state?.from;

  if (!from) {
    return appRoutes.dashboard;
  }

  return `${from.pathname}${from.search}${from.hash}`;
};
