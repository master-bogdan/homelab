import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';
import { Box, Link, Stack, Typography } from '@mui/material';
import type { PropsWithChildren } from 'react';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { OverlineText } from '@/shared/ui';

export interface AuthShellProps extends PropsWithChildren {
  readonly pattern?: 'ambient' | 'dots';
}

const footerLinks = ['Privacy', 'Terms', 'Security'] as const;

export const AuthShell = ({
  children,
  pattern = 'ambient'
}: AuthShellProps) => (
  <Box
    sx={{
      backgroundColor: 'background.default',
      display: 'flex',
      flexDirection: 'column',
      minHeight: '100vh',
      overflow: 'hidden',
      position: 'relative'
    }}
  >
    <Box
      aria-hidden
      sx={{
        inset: 0,
        pointerEvents: 'none',
        position: 'absolute',
        backgroundImage:
          pattern === 'dots'
            ? (theme) =>
                `radial-gradient(${theme.app.borders.ghost} 0.6px, transparent 0.6px), linear-gradient(90deg, rgba(81, 72, 215, 0.04) 0%, transparent 36%, rgba(81, 72, 215, 0.05) 72%, transparent 100%)`
            : 'linear-gradient(90deg, rgba(81, 72, 215, 0.04) 0%, transparent 36%, rgba(81, 72, 215, 0.05) 72%, transparent 100%)',
        backgroundPosition: pattern === 'dots' ? '0 0, 0 0' : '0 0',
        backgroundSize: pattern === 'dots' ? '24px 24px, auto' : 'auto',
        opacity: pattern === 'dots' ? 0.8 : 1
      }}
    />
    <Box
      aria-hidden
      sx={{
        bgcolor: 'primary.main',
        borderRadius: '50%',
        filter: 'blur(120px)',
        height: 320,
        opacity: 0.08,
        position: 'absolute',
        right: '-10%',
        top: '8%',
        width: 360
      }}
    />
    <Box
      aria-hidden
      sx={{
        bgcolor: 'secondary.main',
        borderRadius: '50%',
        bottom: '-10%',
        filter: 'blur(120px)',
        height: 260,
        left: '-8%',
        opacity: 0.22,
        position: 'absolute',
        width: 300
      }}
    />

    <Box
      component="header"
      sx={{
        backdropFilter: (theme) => `blur(${theme.app.effects.backdropBlur})`,
        backgroundColor: (theme) => theme.app.surfaces.overlay,
        left: 0,
        position: 'sticky',
        px: { xs: 2.5, md: 3 },
        py: 2,
        top: 0,
        zIndex: 1
      }}
    >
      <Link
        color="inherit"
        component={RouterLink}
        sx={{
          alignItems: 'center',
          display: 'inline-flex',
          gap: 1,
          textDecoration: 'none'
        }}
        to={appRoutes.login}
        underline="none"
      >
        <ArchitectureRoundedIcon color="primary" />
        <Typography color="text.primary" variant="h6">
          EstimateRoom
        </Typography>
      </Link>
    </Box>

    <Box
      component="main"
      sx={{
        alignItems: 'center',
        display: 'flex',
        flex: 1,
        justifyContent: 'center',
        px: 2,
        position: 'relative',
        py: { xs: 6, md: 8 },
        zIndex: 1
      }}
    >
      <Box sx={{ width: '100%', maxWidth: 520 }}>{children}</Box>
    </Box>

    <Box
      component="footer"
      sx={{
        borderTop: (theme) => `1px solid ${theme.app.borders.ghost}`,
        position: 'relative',
        zIndex: 1
      }}
    >
      <Stack
        alignItems={{ xs: 'flex-start', md: 'center' }}
        direction={{ xs: 'column', md: 'row' }}
        justifyContent="space-between"
        spacing={2}
        sx={{ px: { xs: 2.5, md: 3 }, py: 2.5 }}
      >
        <OverlineText>© 2026 EstimateRoom. All rights reserved.</OverlineText>
        <Stack direction="row" spacing={3}>
          {footerLinks.map((label) => (
            <Link
              key={label}
              color="text.secondary"
              href="#"
              onClick={(event) => event.preventDefault()}
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
