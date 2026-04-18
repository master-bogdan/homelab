import ArrowBackRoundedIcon from '@mui/icons-material/ArrowBackRounded';
import { Link as RouterLink } from 'react-router-dom';

import { AppRoutes } from '@/app/router/routePaths';
import { AppLink } from '@/shared/components';

import {
  authBackToSignInCenteredLinkSx,
  authBackToSignInFormLinkSx,
  authBackToSignInLinkSx
} from './styles';

export const AuthBackToSignInLinkPlacements = {
  CENTERED: 'centered',
  DEFAULT: 'default',
  FORM: 'form'
} as const;

type AuthBackToSignInLinkPlacement =
  (typeof AuthBackToSignInLinkPlacements)[keyof typeof AuthBackToSignInLinkPlacements];

interface AuthBackToSignInLinkProps {
  readonly color?: 'primary' | 'text.secondary';
  readonly placement?: AuthBackToSignInLinkPlacement;
  readonly variant?: 'body2' | 'overline';
}

const getBackLinkSx = (placement: AuthBackToSignInLinkProps['placement']) => {
  if (placement === AuthBackToSignInLinkPlacements.CENTERED) {
    return authBackToSignInCenteredLinkSx;
  }

  if (placement === AuthBackToSignInLinkPlacements.FORM) {
    return authBackToSignInFormLinkSx;
  }

  return authBackToSignInLinkSx;
};

export const AuthBackToSignInLink = ({
  color = 'primary',
  placement = AuthBackToSignInLinkPlacements.DEFAULT,
  variant = 'overline'
}: AuthBackToSignInLinkProps) => (
  <AppLink
    color={color}
    component={RouterLink}
    sx={getBackLinkSx(placement)}
    to={AppRoutes.LOGIN}
    underline="none"
    variant={variant}
  >
    <ArrowBackRoundedIcon fontSize="inherit" />
    Back to Sign In
  </AppLink>
);
