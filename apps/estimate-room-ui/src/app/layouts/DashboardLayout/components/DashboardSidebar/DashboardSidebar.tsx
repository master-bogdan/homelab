import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';

import {
  AppBox,
  AppDrawer,
  AppList,
  AppStack,
  AppTypography
} from '@/shared/components';

import {
  primaryNavigationItems,
  secondaryNavigationItems
} from '../../constants';
import {
  dashboardSidebarBrandIconSx,
  dashboardSidebarBrandSx,
  dashboardSidebarContentSx,
  dashboardSidebarDrawerSx,
  dashboardSidebarListSx,
  dashboardSidebarSecondaryListSx
} from './DashboardSidebar.styles';
import { DashboardSidebarItem } from './DashboardSidebarItem';

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
}: DashboardSidebarProps) => (
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
        {primaryNavigationItems.map((item) => (
          <DashboardSidebarItem
            isDesktop={isDesktop}
            item={item}
            key={item.to}
            onClose={onClose}
            pathname={pathname}
          />
        ))}
      </AppList>
      <AppList disablePadding sx={dashboardSidebarSecondaryListSx}>
        {secondaryNavigationItems.map((item) => (
          <DashboardSidebarItem
            isDesktop={isDesktop}
            item={item}
            key={item.to}
            onClose={onClose}
            pathname={pathname}
          />
        ))}
      </AppList>
    </AppBox>
  </AppDrawer>
);
