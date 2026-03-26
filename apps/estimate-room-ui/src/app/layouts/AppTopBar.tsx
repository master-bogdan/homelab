import MenuRoundedIcon from '@mui/icons-material/MenuRounded';
import LogoutRoundedIcon from '@mui/icons-material/LogoutRounded';
import DarkModeRoundedIcon from '@mui/icons-material/DarkModeRounded';
import LightModeRoundedIcon from '@mui/icons-material/LightModeRounded';
import {
  AppBar,
  Avatar,
  Box,
  IconButton,
  Stack,
  Toolbar,
  Typography
} from '@mui/material';

import { useAppDispatch, useAppSelector } from '@/app/store/hooks';
import { toggleThemeMode } from '@/app/store/uiSlice';
import { appConfig } from '@/shared/config/env';
import { AppButton } from '@/shared/ui';
import { useLogout } from '@/modules/auth/hooks';
import { selectAuthUser } from '@/modules/auth/selectors';
import { selectThemeMode } from '@/app/store/uiSelectors';

export interface AppTopBarProps {
  readonly drawerWidth: number;
  readonly isDesktop: boolean;
  readonly onMenuClick: () => void;
}

export const AppTopBar = ({
  drawerWidth,
  isDesktop,
  onMenuClick
}: AppTopBarProps) => {
  const dispatch = useAppDispatch();
  const { isLoggingOut, logout } = useLogout();
  const user = useAppSelector(selectAuthUser);
  const themeMode = useAppSelector(selectThemeMode);
  const displayName = user?.displayName ?? 'Estimate Room';

  const initials = displayName
    .split(' ')
    .map((namePart) => namePart[0])
    .join('')
    .slice(0, 2)
    .toUpperCase();

  return (
    <AppBar
      color="inherit"
      elevation={0}
      position="fixed"
      sx={{
        backdropFilter: (theme) => `blur(${theme.app.effects.backdropBlur})`,
        backgroundColor: (theme) => theme.app.surfaces.overlay,
        boxShadow: 'none',
        width: { lg: `calc(100% - ${drawerWidth}px)` },
        ml: { lg: `${drawerWidth}px` }
      }}
    >
      <Toolbar sx={{ gap: 2, minHeight: 76, px: { xs: 2, md: 3 } }}>
        {!isDesktop ? (
          <IconButton color="inherit" edge="start" onClick={onMenuClick}>
            <MenuRoundedIcon />
          </IconButton>
        ) : null}
        <Box sx={{ flexGrow: 1, minWidth: 0 }}>
          <Typography color="text.secondary" noWrap variant="overline">
            Workspace
          </Typography>
          <Typography noWrap variant="h6">
            {appConfig.appName}
          </Typography>
          <Typography color="text.secondary" noWrap variant="body2">
            Backend-ready workspace for estimate operations.
          </Typography>
        </Box>
        <Stack alignItems="center" direction="row" spacing={1.5}>
          <IconButton
            color="inherit"
            onClick={() => dispatch(toggleThemeMode())}
            sx={{
              bgcolor: (theme) =>
                theme.palette.mode === 'light'
                  ? theme.app.surfaces.inset
                  : theme.palette.action.hover
            }}
            title="Toggle theme"
          >
            {themeMode === 'light' ? <DarkModeRoundedIcon /> : <LightModeRoundedIcon />}
          </IconButton>
          <Stack alignItems="center" direction="row" spacing={1}>
            <Avatar
              sx={{
                backgroundImage: (theme) => theme.app.gradients.primary,
                color: 'primary.contrastText'
              }}
            >
              {initials}
            </Avatar>
            <Box sx={{ display: { xs: 'none', sm: 'block' } }}>
              <Typography variant="body2">
                {user?.displayName ?? 'Pending backend auth'}
              </Typography>
              <Typography color="text.secondary" variant="caption">
                {user?.email ?? 'No active session'}
              </Typography>
            </Box>
          </Stack>
          <AppButton
            color="secondary"
            disabled={isLoggingOut}
            loading={isLoggingOut}
            loadingText="Logging Out..."
            onClick={() => {
              void logout();
            }}
            startIcon={<LogoutRoundedIcon />}
            variant="contained"
          >
            Logout
          </AppButton>
        </Stack>
      </Toolbar>
    </AppBar>
  );
};
