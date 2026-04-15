import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';
import {
  Box,
  Drawer,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Stack,
  Typography
} from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

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
        <ListItemButton
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
          <ListItemIcon sx={dashboardSidebarItemIconSx}>
            <Icon fontSize="small" />
          </ListItemIcon>
          <ListItemText
            primary={label}
            primaryTypographyProps={{
              fontWeight: isSelected ? 700 : 600,
              variant: 'body2'
            }}
          />
        </ListItemButton>
      );
    });

  return (
    <Drawer
      ModalProps={{ keepMounted: true }}
      onClose={onClose}
      open={isDesktop ? true : isOpen}
      sx={dashboardSidebarDrawerSx}
      variant={isDesktop ? 'permanent' : 'temporary'}
    >
      <Box sx={dashboardSidebarContentSx}>
        <Stack direction="row" spacing={1.5} sx={dashboardSidebarBrandSx}>
          <Box sx={dashboardSidebarBrandIconSx}>
            <ArchitectureRoundedIcon fontSize="small" />
          </Box>
          <Stack spacing={0.25}>
            <Typography variant="h6">EstimateRoom</Typography>
            <Typography color="primary.main" variant="overline">
              {occupationLabel}
            </Typography>
          </Stack>
        </Stack>
        <List disablePadding sx={dashboardSidebarListSx}>
          {renderNavigationItems(primaryNavigationItems)}
        </List>
        <List disablePadding sx={dashboardSidebarSecondaryListSx}>
          {renderNavigationItems(secondaryNavigationItems)}
        </List>
      </Box>
    </Drawer>
  );
};
