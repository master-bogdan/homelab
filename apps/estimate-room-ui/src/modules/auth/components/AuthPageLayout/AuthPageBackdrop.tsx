import { AppBox } from '@/shared/components';

import {
  type AuthPageLayoutPattern,
  authPageLayoutGlowBottomSx,
  authPageLayoutGlowTopSx,
  getAuthPageLayoutBackdropSx
} from './styles';

interface AuthPageBackdropProps {
  readonly pattern: AuthPageLayoutPattern;
}

export const AuthPageBackdrop = ({ pattern }: AuthPageBackdropProps) => (
  <>
    <AppBox aria-hidden sx={getAuthPageLayoutBackdropSx(pattern)} />
    <AppBox aria-hidden sx={authPageLayoutGlowTopSx} />
    <AppBox aria-hidden sx={authPageLayoutGlowBottomSx} />
  </>
);
