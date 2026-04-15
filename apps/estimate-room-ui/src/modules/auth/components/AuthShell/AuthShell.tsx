import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';
import type { PropsWithChildren } from 'react';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppBox, AppLink, AppStack, AppTypography, OverlineText } from '@/shared/ui';

import {
  type AuthShellPattern,
  authShellFooterRootSx,
  authShellFooterStackSx,
  authShellGlowBottomSx,
  authShellGlowTopSx,
  authShellHeaderRootSx,
  authShellHomeLinkSx,
  authShellInnerSx,
  authShellMainRootSx,
  authShellRootSx,
  authShellUtilityLinkSx,
  getAuthShellBackdropSx
} from './styles';

export interface AuthShellProps extends PropsWithChildren {
  readonly pattern?: AuthShellPattern;
}

const footerLinks = ['Privacy', 'Terms', 'Security'] as const;

export const AuthShell = ({
  children,
  pattern = 'ambient'
}: AuthShellProps) => (
  <AppBox sx={authShellRootSx}>
    <AppBox aria-hidden sx={getAuthShellBackdropSx(pattern)} />
    <AppBox aria-hidden sx={authShellGlowTopSx} />
    <AppBox aria-hidden sx={authShellGlowBottomSx} />

    <AppBox component="header" sx={authShellHeaderRootSx}>
      <AppLink
        color="inherit"
        component={RouterLink}
        sx={authShellHomeLinkSx}
        to={appRoutes.login}
        underline="none"
      >
        <ArchitectureRoundedIcon color="primary" />
        <AppTypography color="text.primary" variant="h6">
          EstimateRoom
        </AppTypography>
      </AppLink>
    </AppBox>

    <AppBox component="main" sx={authShellMainRootSx}>
      <AppBox sx={authShellInnerSx}>{children}</AppBox>
    </AppBox>

    <AppBox component="footer" sx={authShellFooterRootSx}>
      <AppStack
        alignItems={{ xs: 'flex-start', md: 'center' }}
        direction={{ xs: 'column', md: 'row' }}
        justifyContent="space-between"
        spacing={2}
        sx={authShellFooterStackSx}
      >
        <OverlineText>© 2026 EstimateRoom. All rights reserved.</OverlineText>
        <AppStack direction="row" spacing={3}>
          {footerLinks.map((label) => (
            <AppLink
              key={label}
              color="text.secondary"
              href="#"
              onClick={(event) => event.preventDefault()}
              sx={authShellUtilityLinkSx}
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
