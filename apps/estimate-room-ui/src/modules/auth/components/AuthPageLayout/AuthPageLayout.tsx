import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';
import type { PropsWithChildren } from 'react';
import { Link as RouterLink } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';
import { AppBox, AppLink, AppStack, AppTypography, OverlineText } from '@/shared/ui';

import {
  type AuthPageLayoutPattern,
  authPageLayoutFooterRootSx,
  authPageLayoutFooterStackSx,
  authPageLayoutGlowBottomSx,
  authPageLayoutGlowTopSx,
  authPageLayoutHeaderRootSx,
  authPageLayoutHomeLinkSx,
  authPageLayoutInnerSx,
  authPageLayoutMainRootSx,
  authPageLayoutRootSx,
  authPageLayoutUtilityLinkSx,
  getAuthPageLayoutBackdropSx
} from './styles';

export interface AuthPageLayoutProps extends PropsWithChildren {
  readonly pattern?: AuthPageLayoutPattern;
}

const footerLinks = ['Privacy', 'Terms', 'Security'] as const;

export const AuthPageLayout = ({
  children,
  pattern = 'ambient'
}: AuthPageLayoutProps) => (
  <AppBox sx={authPageLayoutRootSx}>
    <AppBox aria-hidden sx={getAuthPageLayoutBackdropSx(pattern)} />
    <AppBox aria-hidden sx={authPageLayoutGlowTopSx} />
    <AppBox aria-hidden sx={authPageLayoutGlowBottomSx} />

    <AppBox component="header" sx={authPageLayoutHeaderRootSx}>
      <AppLink
        color="inherit"
        component={RouterLink}
        sx={authPageLayoutHomeLinkSx}
        to={AppRoutes.LOGIN}
        underline="none"
      >
        <ArchitectureRoundedIcon color="primary" />
        <AppTypography color="text.primary" variant="h6">
          EstimateRoom
        </AppTypography>
      </AppLink>
    </AppBox>

    <AppBox component="main" sx={authPageLayoutMainRootSx}>
      <AppBox sx={authPageLayoutInnerSx}>{children}</AppBox>
    </AppBox>

    <AppBox component="footer" sx={authPageLayoutFooterRootSx}>
      <AppStack
        alignItems={{ xs: 'flex-start', md: 'center' }}
        direction={{ xs: 'column', md: 'row' }}
        justifyContent="space-between"
        spacing={2}
        sx={authPageLayoutFooterStackSx}
      >
        <OverlineText>© 2026 EstimateRoom. All rights reserved.</OverlineText>
        <AppStack direction="row" spacing={3}>
          {footerLinks.map((label) => (
            <AppLink
              key={label}
              color="text.secondary"
              href="#"
              onClick={(event) => event.preventDefault()}
              sx={authPageLayoutUtilityLinkSx}
              underline="always"
              variant="overline"
            >
              {label}
            </AppLink>
          ))}
        </AppStack>
      </AppStack>
    </AppBox>
  </AppBox>
);
