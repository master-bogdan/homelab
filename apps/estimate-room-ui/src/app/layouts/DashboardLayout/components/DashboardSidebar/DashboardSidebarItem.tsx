import { Link as RouterLink } from 'react-router-dom';

import {
  AppListItemButton,
  AppListItemIcon,
  AppListItemText,
  AppTypography
} from '@/shared/components';

import type { DashboardLayoutNavigationItem } from '../../constants';
import {
  dashboardSidebarItemIconSx,
  dashboardSidebarItemSx,
  getDashboardSidebarItemLabelSx
} from './DashboardSidebar.styles';
import { isPathSelected } from '../../utils';

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
  const handleClick = isDesktop ? undefined : onClose;

  return (
    <AppListItemButton
      component={RouterLink}
      onClick={handleClick}
      selected={isSelected}
      sx={dashboardSidebarItemSx}
      to={to}
    >
      <AppListItemIcon sx={dashboardSidebarItemIconSx}>
        <Icon fontSize="small" />
      </AppListItemIcon>
      <AppListItemText
        disableTypography
        primary={
          <AppTypography
            component="span"
            sx={getDashboardSidebarItemLabelSx(isSelected)}
            variant="body2"
          >
            {label}
          </AppTypography>
        }
      />
    </AppListItemButton>
  );
};
