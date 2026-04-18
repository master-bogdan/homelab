import { AppRoutes } from '@/app/router/routePaths';

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
    return AppRoutes.DASHBOARD;
  }

  return `${from.pathname}${from.search}${from.hash}`;
};
