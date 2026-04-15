import MenuRoundedIcon from '@mui/icons-material/MenuRounded';
import NotificationsRoundedIcon from '@mui/icons-material/NotificationsRounded';
import { Avatar, Box, IconButton, Stack, Toolbar, Typography } from '@mui/material';
import type { MouseEvent } from 'react';

import { AppButton } from '@/shared/ui';

import type { DashboardLayoutMeta } from '../dashboardLayout.meta';
import {
  dashboardHeaderActionsSx,
  dashboardHeaderAvatarButtonSx,
  dashboardHeaderAvatarSx,
  dashboardHeaderDividerSx,
  dashboardHeaderMenuButtonSx,
  dashboardHeaderRootSx,
  dashboardHeaderTitleSx,
  dashboardHeaderToolbarSx
} from './DashboardHeader.styles';

export interface DashboardHeaderProps {
  readonly isDesktop: boolean;
  readonly onOpenCreateRoom: () => void;
  readonly onOpenJoinRoom: () => void;
  readonly onOpenSidebar: () => void;
  readonly onOpenUserMenu: (anchor: HTMLElement) => void;
  readonly routeMeta: DashboardLayoutMeta;
  readonly userAvatarUrl: string | null;
  readonly userInitials: string;
}

export const DashboardHeader = ({
  isDesktop,
  onOpenCreateRoom,
  onOpenJoinRoom,
  onOpenSidebar,
  onOpenUserMenu,
  routeMeta,
  userAvatarUrl,
  userInitials
}: DashboardHeaderProps) => (
  <Box component="header" sx={dashboardHeaderRootSx}>
    <Toolbar sx={dashboardHeaderToolbarSx}>
      {!isDesktop ? (
        <IconButton
          color="inherit"
          edge="start"
          onClick={onOpenSidebar}
          sx={dashboardHeaderMenuButtonSx}
        >
          <MenuRoundedIcon />
        </IconButton>
      ) : null}
      <Stack minWidth={0} spacing={0.5} sx={dashboardHeaderTitleSx}>
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
        sx={dashboardHeaderActionsSx}
      >
        <AppButton color="secondary" onClick={onOpenJoinRoom} variant="contained">
          Join room
        </AppButton>
        <AppButton onClick={onOpenCreateRoom} variant="contained">
          Create room
        </AppButton>
        <Box sx={dashboardHeaderDividerSx} />
        <IconButton aria-label="Notifications" color="inherit">
          <NotificationsRoundedIcon />
        </IconButton>
        <IconButton
          aria-label="Open user menu"
          color="inherit"
          onClick={(event: MouseEvent<HTMLElement>) => {
            onOpenUserMenu(event.currentTarget);
          }}
          sx={dashboardHeaderAvatarButtonSx}
        >
          <Avatar src={userAvatarUrl ?? undefined} sx={dashboardHeaderAvatarSx}>
            {userInitials}
          </Avatar>
        </IconButton>
      </Stack>
    </Toolbar>
  </Box>
);
