import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';
import DashboardRoundedIcon from '@mui/icons-material/DashboardRounded';
import HistoryRoundedIcon from '@mui/icons-material/HistoryRounded';
import LogoutRoundedIcon from '@mui/icons-material/LogoutRounded';
import MenuRoundedIcon from '@mui/icons-material/MenuRounded';
import NotificationsRoundedIcon from '@mui/icons-material/NotificationsRounded';
import PersonRoundedIcon from '@mui/icons-material/PersonRounded';
import SettingsRoundedIcon from '@mui/icons-material/SettingsRounded';
import {
  Avatar,
  Box,
  Drawer,
  IconButton,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Menu,
  MenuItem,
  Stack,
  Toolbar,
  Typography,
  useMediaQuery
} from '@mui/material';
import type { MouseEvent } from 'react';
import { useEffect, useState } from 'react';
import { Link as RouterLink, Outlet, useLocation } from 'react-router-dom';

import { useAppDispatch, useAppSelector } from '@/app/store/hooks';
import { closeSidebar, setSidebarOpen } from '@/app/store/uiSlice';
import { selectIsSidebarOpen } from '@/app/store/uiSelectors';
import { useLogout } from '@/modules/auth/hooks';
import { selectAuthUser } from '@/modules/auth/store';
import { DashboardDialogs, useDashboardActions } from '@/modules/dashboard';
import { APP_DRAWER_WIDTH } from '@/shared/constants/layout';
import { appRoutes } from '@/shared/constants/routes';
import { AppButton } from '@/shared/ui';
import { getInitials } from '@/modules/dashboard/utils';

import { resolveDashboardLayoutMeta } from './dashboardLayoutMeta';

const primaryNavigationItems = [
  {
    icon: DashboardRoundedIcon,
    label: 'Dashboard',
    to: appRoutes.dashboard
  },
  {
    icon: HistoryRoundedIcon,
    label: 'History',
    to: appRoutes.history
  }
] as const;

const secondaryNavigationItems = [
  {
    icon: PersonRoundedIcon,
    label: 'Profile',
    to: appRoutes.profile
  },
  {
    icon: SettingsRoundedIcon,
    label: 'Settings',
    to: appRoutes.settings
  }
] as const;

const isPathSelected = (pathname: string, itemPath: string) =>
  pathname === itemPath || pathname.startsWith(`${itemPath}/`);

const DashboardLayoutContent = () => {
  const dispatch = useAppDispatch();
  const location = useLocation();
  const isDesktop = useMediaQuery((theme) => theme.breakpoints.up('lg'));
  const isSidebarOpen = useAppSelector(selectIsSidebarOpen);
  const user = useAppSelector(selectAuthUser);
  const { isLoggingOut, logout } = useLogout();
  const { openCreateRoom, openJoinRoom } = useDashboardActions();
  const [userMenuAnchor, setUserMenuAnchor] = useState<HTMLElement | null>(null);
  const routeMeta = resolveDashboardLayoutMeta(location.pathname);
  const displayName = user?.displayName ?? 'EstimateRoom User';
  const occupationLabel = user?.occupation?.trim() || 'EstimateRoom Member';
  const userInitials = getInitials(displayName);

  useEffect(() => {
    dispatch(setSidebarOpen(isDesktop));
  }, [dispatch, isDesktop]);

  const navigationContent = (
    <Box
      sx={{
        bgcolor: (theme) => theme.app.surfaces.section,
        display: 'flex',
        flexDirection: 'column',
        height: '100%',
        px: 2,
        py: 3
      }}
    >
      <Stack direction="row" spacing={1.5} sx={{ mb: 5, px: 1 }}>
        <Box
          sx={{
            alignItems: 'center',
            backgroundImage: (theme) => theme.app.gradients.primary,
            borderRadius: 1.5,
            color: 'primary.contrastText',
            display: 'flex',
            height: 40,
            justifyContent: 'center',
            width: 40
          }}
        >
          <ArchitectureRoundedIcon fontSize="small" />
        </Box>
        <Stack spacing={0.25}>
          <Typography variant="h6">EstimateRoom</Typography>
          <Typography color="primary.main" variant="overline">
            {occupationLabel}
          </Typography>
        </Stack>
      </Stack>
      <List disablePadding sx={{ display: 'grid', gap: 0.5 }}>
        {primaryNavigationItems.map(({ icon: Icon, label, to }) => (
          <ListItemButton
            component={RouterLink}
            key={to}
            onClick={() => {
              if (!isDesktop) {
                dispatch(closeSidebar());
              }
            }}
            selected={isPathSelected(location.pathname, to)}
            sx={{ borderRadius: 1.5 }}
            to={to}
          >
            <ListItemIcon sx={{ color: 'inherit', minWidth: 40 }}>
              <Icon fontSize="small" />
            </ListItemIcon>
            <ListItemText
              primary={label}
              primaryTypographyProps={{
                fontWeight: isPathSelected(location.pathname, to) ? 700 : 600,
                variant: 'body2'
              }}
            />
          </ListItemButton>
        ))}
      </List>
      <List disablePadding sx={{ display: 'grid', gap: 0.5, mt: 'auto' }}>
        {secondaryNavigationItems.map(({ icon: Icon, label, to }) => (
          <ListItemButton
            component={RouterLink}
            key={to}
            onClick={() => {
              if (!isDesktop) {
                dispatch(closeSidebar());
              }
            }}
            selected={isPathSelected(location.pathname, to)}
            sx={{ borderRadius: 1.5 }}
            to={to}
          >
            <ListItemIcon sx={{ color: 'inherit', minWidth: 40 }}>
              <Icon fontSize="small" />
            </ListItemIcon>
            <ListItemText
              primary={label}
              primaryTypographyProps={{
                fontWeight: isPathSelected(location.pathname, to) ? 700 : 600,
                variant: 'body2'
              }}
            />
          </ListItemButton>
        ))}
      </List>
    </Box>
  );

  return (
    <Box sx={{ bgcolor: 'background.default', display: 'flex', minHeight: '100vh' }}>
      <Drawer
        ModalProps={{ keepMounted: true }}
        onClose={() => dispatch(closeSidebar())}
        open={isDesktop ? true : isSidebarOpen}
        sx={{
          '& .MuiDrawer-paper': {
            borderRight: 'none',
            boxSizing: 'border-box',
            width: APP_DRAWER_WIDTH
          },
          width: APP_DRAWER_WIDTH
        }}
        variant={isDesktop ? 'permanent' : 'temporary'}
      >
        {navigationContent}
      </Drawer>
      <Box
        sx={{
          display: 'flex',
          flex: 1,
          flexDirection: 'column',
          minWidth: 0
        }}
      >
        <Box
          component="header"
          sx={{
            backdropFilter: (theme) => `blur(${theme.app.effects.backdropBlur})`,
            bgcolor: (theme) => theme.app.surfaces.overlay,
            borderBottom: (theme) => `1px solid ${theme.app.borders.ghost}`,
            position: 'sticky',
            top: 0,
            zIndex: (theme) => theme.zIndex.appBar
          }}
        >
          <Toolbar
            sx={{
              alignItems: { xs: 'flex-start', md: 'center' },
              flexWrap: 'wrap',
              gap: 2,
              minHeight: '72px !important',
              px: { xs: 2, md: 4 },
              py: 1.5
            }}
          >
            {!isDesktop ? (
              <IconButton
                color="inherit"
                edge="start"
                onClick={() => dispatch(setSidebarOpen(true))}
                sx={{ mt: { xs: 0.5, md: 0 } }}
              >
                <MenuRoundedIcon />
              </IconButton>
            ) : null}
            <Stack minWidth={0} spacing={0.5} sx={{ flex: 1 }}>
              <Typography variant="h6">{routeMeta.title}</Typography>
              <Typography color="text.secondary" variant="body2">
                {routeMeta.description}
              </Typography>
            </Stack>
            <Stack
              alignItems="center"
              direction="row"
              flexWrap="wrap"
              gap={1.25}
              justifyContent="flex-end"
              sx={{ ml: 'auto' }}
            >
              <AppButton color="secondary" onClick={openJoinRoom} variant="contained">
                Join room
              </AppButton>
              <AppButton onClick={openCreateRoom} variant="contained">
                Create room
              </AppButton>
              <Box
                sx={{
                  alignSelf: 'stretch',
                  bgcolor: (theme) => theme.app.borders.ghost,
                  display: { xs: 'none', sm: 'block' },
                  width: '1px'
                }}
              />
              <IconButton aria-label="Notifications" color="inherit">
                <NotificationsRoundedIcon />
              </IconButton>
              <IconButton
                aria-label="Open user menu"
                color="inherit"
                onClick={(event: MouseEvent<HTMLElement>) => {
                  setUserMenuAnchor(event.currentTarget);
                }}
                sx={{ p: 0 }}
              >
                <Avatar
                  src={user?.avatarUrl ?? undefined}
                  sx={{
                    bgcolor: 'primary.main',
                    color: 'primary.contrastText',
                    height: 36,
                    width: 36
                  }}
                >
                  {userInitials}
                </Avatar>
              </IconButton>
            </Stack>
          </Toolbar>
        </Box>
        <Box
          component="main"
          sx={{
            flex: 1,
            px: { xs: 2, md: 4 },
            py: { xs: 3, md: 4 }
          }}
        >
          <Outlet />
        </Box>
      </Box>
      <Menu
        anchorEl={userMenuAnchor}
        anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
        onClose={() => setUserMenuAnchor(null)}
        open={Boolean(userMenuAnchor)}
        PaperProps={{
          elevation: 0,
          sx: {
            border: (theme) => `1px solid ${theme.app.borders.ghost}`,
            boxShadow: (theme) => theme.app.effects.ambientShadow,
            mt: 1,
            minWidth: 240
          }
        }}
        transformOrigin={{ horizontal: 'right', vertical: 'top' }}
      >
        <Stack spacing={0.25} sx={{ px: 2, py: 1.5 }}>
          <Typography variant="subtitle2">{displayName}</Typography>
          <Typography color="text.secondary" variant="caption">
            {user?.email ?? 'No active session'}
          </Typography>
        </Stack>
        <MenuItem
          component={RouterLink}
          onClick={() => setUserMenuAnchor(null)}
          to={appRoutes.profile}
        >
          <ListItemIcon>
            <PersonRoundedIcon fontSize="small" />
          </ListItemIcon>
          Profile
        </MenuItem>
        <MenuItem
          component={RouterLink}
          onClick={() => setUserMenuAnchor(null)}
          to={appRoutes.settings}
        >
          <ListItemIcon>
            <SettingsRoundedIcon fontSize="small" />
          </ListItemIcon>
          Settings
        </MenuItem>
        <MenuItem
          disabled={isLoggingOut}
          onClick={() => {
            setUserMenuAnchor(null);
            void logout();
          }}
        >
          <ListItemIcon>
            <LogoutRoundedIcon fontSize="small" />
          </ListItemIcon>
          {isLoggingOut ? 'Logging out...' : 'Log out'}
        </MenuItem>
      </Menu>
      <DashboardDialogs />
    </Box>
  );
};

export const DashboardLayout = () => <DashboardLayoutContent />;
