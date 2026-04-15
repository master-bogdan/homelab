import { Box, useMediaQuery } from '@mui/material';
import { useEffect, useState } from 'react';
import { Outlet, useLocation } from 'react-router-dom';

import { useAppDispatch, useAppSelector } from '@/shared/store';
import { closeSidebar, setSidebarOpen } from '@/modules/system';
import { selectIsSidebarOpen } from '@/modules/system';
import { selectAuthUser, useLogout } from '@/modules/auth';
import { DashboardDialogs, getInitials, useDashboardActions } from '@/modules/dashboard';

import { resolveDashboardLayoutMeta } from '../dashboardLayout.meta';
import { DashboardHeader } from '../DashboardHeader';
import { DashboardSidebar } from '../DashboardSidebar';
import { DashboardUserMenu } from '../DashboardUserMenu';
import {
  dashboardLayoutBodySx,
  dashboardLayoutMainSx,
  dashboardLayoutRootSx
} from './DashboardLayoutContent.styles';

export const DashboardLayoutContent = () => {
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

  return (
    <Box sx={dashboardLayoutRootSx}>
      <DashboardSidebar
        isDesktop={isDesktop}
        isOpen={isSidebarOpen}
        occupationLabel={occupationLabel}
        onClose={() => dispatch(closeSidebar())}
        pathname={location.pathname}
      />
      <Box sx={dashboardLayoutBodySx}>
        <DashboardHeader
          isDesktop={isDesktop}
          onOpenCreateRoom={openCreateRoom}
          onOpenJoinRoom={openJoinRoom}
          onOpenSidebar={() => dispatch(setSidebarOpen(true))}
          onOpenUserMenu={setUserMenuAnchor}
          routeMeta={routeMeta}
          userAvatarUrl={user?.avatarUrl ?? null}
          userInitials={userInitials}
        />
        <Box component="main" sx={dashboardLayoutMainSx}>
          <Outlet />
        </Box>
      </Box>
      <DashboardUserMenu
        anchorEl={userMenuAnchor}
        displayName={displayName}
        email={user?.email ?? null}
        isLoggingOut={isLoggingOut}
        onClose={() => setUserMenuAnchor(null)}
        onLogout={logout}
      />
      <DashboardDialogs />
    </Box>
  );
};
