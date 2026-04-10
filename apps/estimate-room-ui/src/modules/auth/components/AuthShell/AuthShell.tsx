import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';
import { Box, Link, Stack, Typography } from '@mui/material';
import type { PropsWithChildren } from 'react';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { OverlineText } from '@/shared/ui';

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
  <Box sx={authShellRootSx}>
    <Box aria-hidden sx={getAuthShellBackdropSx(pattern)} />
    <Box aria-hidden sx={authShellGlowTopSx} />
    <Box aria-hidden sx={authShellGlowBottomSx} />

    <Box component="header" sx={authShellHeaderRootSx}>
      <Link
        color="inherit"
        component={RouterLink}
        sx={authShellHomeLinkSx}
        to={appRoutes.login}
        underline="none"
      >
        <ArchitectureRoundedIcon color="primary" />
        <Typography color="text.primary" variant="h6">
          EstimateRoom
        </Typography>
      </Link>
    </Box>

    <Box component="main" sx={authShellMainRootSx}>
      <Box sx={authShellInnerSx}>{children}</Box>
    </Box>

    <Box component="footer" sx={authShellFooterRootSx}>
      <Stack
        alignItems={{ xs: 'flex-start', md: 'center' }}
        direction={{ xs: 'column', md: 'row' }}
        justifyContent="space-between"
        spacing={2}
        sx={authShellFooterStackSx}
      >
        <OverlineText>© 2026 EstimateRoom. All rights reserved.</OverlineText>
        <Stack direction="row" spacing={3}>
          {footerLinks.map((label) => (
            <Link
              key={label}
              color="text.secondary"
              href="#"
              onClick={(event) => event.preventDefault()}
              sx={authShellUtilityLinkSx}
              underline="always"
              variant="overline"
            >
              {label}
            </Link>
          ))}
        </Stack>
      </Stack>
    </Box>
  </Box>
);
