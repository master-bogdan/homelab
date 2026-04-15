import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';
import { Link as RouterLink } from 'react-router-dom';

import {
  AppBox,
  AppDrawer,
  AppList,
  AppListItemButton,
  AppListItemIcon,
  AppListItemText,
  AppStack,
  AppTypography
} from '@/shared/ui';

import {
  primaryNavigationItems,
  secondaryNavigationItems
} from './dashboardSidebar.constants';
import {
  dashboardSidebarBrandIconSx,
  dashboardSidebarBrandSx,
  dashboardSidebarContentSx,
  dashboardSidebarDrawerSx,
  dashboardSidebarItemIconSx,
  dashboardSidebarItemSx,
  dashboardSidebarListSx,
  dashboardSidebarSecondaryListSx
} from './DashboardSidebar.styles';
import { isPathSelected } from './dashboardSidebar.utils';

export interface DashboardSidebarProps {
  readonly isDesktop: boolean;
  readonly isOpen: boolean;
  readonly occupationLabel: string;
  readonly onClose: () => void;
  readonly pathname: string;
}

export const DashboardSidebar = ({
  isDesktop,
  isOpen,
  occupationLabel,
  onClose,
  pathname
}: DashboardSidebarProps) => {
  const renderNavigationItems = (
    items: typeof primaryNavigationItems | typeof secondaryNavigationItems
  ) =>
    items.map(({ icon: Icon, label, to }) => {
      const isSelected = isPathSelected(pathname, to);

      return (
        <AppListItemButton
          component={RouterLink}
          key={to}
          onClick={() => {
            if (!isDesktop) {
              onClose();
            }
          }}
          selected={isSelected}
          sx={dashboardSidebarItemSx}
          to={to}
        >
          <AppListItemIcon sx={dashboardSidebarItemIconSx}>
            <Icon fontSize="small" />
          </AppListItemIcon>
          <AppListItemText
            primary={label}
            primaryTypographyProps={{
              fontWeight: isSelected ? 700 : 600,
              variant: 'body2'
            }}
          />
        </AppListItemButton>
      );
    });

  return (
    <AppDrawer
      ModalProps={{ keepMounted: true }}
      onClose={onClose}
      open={isDesktop ? true : isOpen}
      sx={dashboardSidebarDrawerSx}
      variant={isDesktop ? 'permanent' : 'temporary'}
    >
      <AppBox sx={dashboardSidebarContentSx}>
        <AppStack direction="row" spacing={1.5} sx={dashboardSidebarBrandSx}>
          <AppBox sx={dashboardSidebarBrandIconSx}>
            <ArchitectureRoundedIcon fontSize="small" />
          </AppBox>
          <AppStack spacing={0.25}>
            <AppTypography variant="h6">EstimateRoom</AppTypography>
            <AppTypography color="primary.main" variant="overline">
              {occupationLabel}
            </AppTypography>
          </AppStack>
        </AppStack>
        <AppList disablePadding sx={dashboardSidebarListSx}>
          {renderNavigationItems(primaryNavigationItems)}
        </AppList>
        <AppList disablePadding sx={dashboardSidebarSecondaryListSx}>
          {renderNavigationItems(secondaryNavigationItems)}
        </AppList>
      </AppBox>
    </AppDrawer>
  );
};
