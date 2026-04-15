import { Link as RouterLink } from 'react-router-dom';

import {
  AppListItemButton,
  AppListItemIcon,
  AppListItemText
} from '@/shared/ui';

import type { DashboardLayoutNavigationItem } from '../../constants';
import {
  dashboardSidebarItemIconSx,
  dashboardSidebarItemSx
} from './DashboardSidebar.styles';
import { isPathSelected } from './dashboardSidebar.utils';

export interface DashboardSidebarItemProps {
  readonly isDesktop: boolean;
  readonly item: DashboardLayoutNavigationItem;
  readonly onClose: () => void;
  readonly pathname: string;
}

export const DashboardSidebarItem = ({
  isDesktop,
  item,
  onClose,
  pathname
}: DashboardSidebarItemProps) => {
  const { icon: Icon, label, to } = item;
  const isSelected = isPathSelected(pathname, to);

  return (
    <AppListItemButton
      component={RouterLink}
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
};
