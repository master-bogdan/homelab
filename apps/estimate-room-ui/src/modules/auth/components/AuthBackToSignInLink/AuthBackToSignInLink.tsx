import ArrowBackRoundedIcon from '@mui/icons-material/ArrowBackRounded';
import { Link as RouterLink } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';
import { AppLink } from '@/shared/ui';

import {
  authBackToSignInCenteredLinkSx,
  authBackToSignInFormLinkSx,
  authBackToSignInLinkSx
} from './styles';

interface AuthBackToSignInLinkProps {
  readonly color?: 'primary' | 'text.secondary';
  readonly placement?: 'default' | 'centered' | 'form';
  readonly variant?: 'body2' | 'overline';
}

const getBackLinkSx = (placement: AuthBackToSignInLinkProps['placement']) => {
  if (placement === 'centered') {
    return authBackToSignInCenteredLinkSx;
  }

  if (placement === 'form') {
    return authBackToSignInFormLinkSx;
  }

  return authBackToSignInLinkSx;
};

export const AuthBackToSignInLink = ({
  color = 'primary',
  placement = 'default',
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
