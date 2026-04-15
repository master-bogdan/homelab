import { useEffect, useState } from 'react';
import { Outlet, useLocation } from 'react-router-dom';

import { useAppMediaQuery } from '@/shared/hooks';
import { AppBox } from '@/shared/ui';
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
  const isDesktop = useAppMediaQuery((theme) => theme.breakpoints.up('lg'));
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
    <AppBox sx={dashboardLayoutRootSx}>
      <DashboardSidebar
        isDesktop={isDesktop}
        isOpen={isSidebarOpen}
        occupationLabel={occupationLabel}
        onClose={() => dispatch(closeSidebar())}
        pathname={location.pathname}
      />
      <AppBox sx={dashboardLayoutBodySx}>
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
        <AppBox component="main" sx={dashboardLayoutMainSx}>
          <Outlet />
        </AppBox>
      </AppBox>
      <DashboardUserMenu
        anchorEl={userMenuAnchor}
        displayName={displayName}
        email={user?.email ?? null}
        isLoggingOut={isLoggingOut}
        onClose={() => setUserMenuAnchor(null)}
        onLogout={logout}
      />
      <DashboardDialogs />
    </AppBox>
  );
};
