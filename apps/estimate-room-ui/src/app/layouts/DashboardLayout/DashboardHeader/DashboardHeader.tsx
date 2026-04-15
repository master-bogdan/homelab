import MenuRoundedIcon from '@mui/icons-material/MenuRounded';
import NotificationsRoundedIcon from '@mui/icons-material/NotificationsRounded';
import type { MouseEvent } from 'react';

import {
  AppAvatar,
  AppBox,
  AppButton,
  AppIconButton,
  AppStack,
  AppToolbar,
  AppTypography
} from '@/shared/ui';

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
  <AppBox component="header" sx={dashboardHeaderRootSx}>
    <AppToolbar sx={dashboardHeaderToolbarSx}>
      {!isDesktop ? (
        <AppIconButton
          color="inherit"
          edge="start"
          onClick={onOpenSidebar}
          sx={dashboardHeaderMenuButtonSx}
        >
          <MenuRoundedIcon />
        </AppIconButton>
      ) : null}
      <AppStack minWidth={0} spacing={0.5} sx={dashboardHeaderTitleSx}>
        <AppTypography variant="h6">{routeMeta.title}</AppTypography>
        <AppTypography color="text.secondary" variant="body2">
          {routeMeta.description}
        </AppTypography>
      </AppStack>
      <AppStack
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
        <AppBox sx={dashboardHeaderDividerSx} />
        <AppIconButton aria-label="Notifications" color="inherit">
          <NotificationsRoundedIcon />
        </AppIconButton>
        <AppIconButton
          aria-label="Open user menu"
          color="inherit"
          onClick={(event: MouseEvent<HTMLElement>) => {
            onOpenUserMenu(event.currentTarget);
          }}
          sx={dashboardHeaderAvatarButtonSx}
        >
          <AppAvatar src={userAvatarUrl ?? undefined} sx={dashboardHeaderAvatarSx}>
            {userInitials}
          </AppAvatar>
        </AppIconButton>
      </AppStack>
    </AppToolbar>
  </AppBox>
);
