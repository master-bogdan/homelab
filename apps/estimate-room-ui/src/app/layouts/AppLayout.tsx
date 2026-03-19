import { useEffect } from 'react';
import { Box, useMediaQuery } from '@mui/material';
import { Outlet } from 'react-router-dom';

import { useAppDispatch } from '@/app/store/hooks';
import { setSidebarOpen } from '@/app/store/uiSlice';
import { APP_DRAWER_WIDTH } from '@/shared/constants/layout';

import { AppSidebar } from './AppSidebar';
import { AppTopBar } from './AppTopBar';
import { ContentShell } from './ContentShell';

export const AppLayout = () => {
  const dispatch = useAppDispatch();
  const isDesktop = useMediaQuery((theme) => theme.breakpoints.up('lg'));

  useEffect(() => {
    dispatch(setSidebarOpen(isDesktop));
  }, [dispatch, isDesktop]);

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh', bgcolor: 'background.default' }}>
      <AppTopBar
        drawerWidth={APP_DRAWER_WIDTH}
        isDesktop={isDesktop}
        onMenuClick={() => dispatch(setSidebarOpen(true))}
      />
      <AppSidebar isDesktop={isDesktop} />
      <ContentShell drawerWidth={APP_DRAWER_WIDTH}>
        <Outlet />
      </ContentShell>
    </Box>
  );
};
