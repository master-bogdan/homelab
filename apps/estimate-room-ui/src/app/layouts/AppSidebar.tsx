import DashboardRoundedIcon from '@mui/icons-material/DashboardRounded';
import HistoryRoundedIcon from '@mui/icons-material/HistoryRounded';
import PersonRoundedIcon from '@mui/icons-material/PersonRounded';
import SettingsRoundedIcon from '@mui/icons-material/SettingsRounded';
import AddHomeWorkRoundedIcon from '@mui/icons-material/AddHomeWorkRounded';
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
import { Link as RouterLink, useLocation } from 'react-router-dom';

import { useAppDispatch, useAppSelector } from '@/app/store/hooks';
import { closeSidebar } from '@/app/store/uiSlice';
import { selectIsSidebarOpen } from '@/app/store/uiSelectors';
import { APP_DRAWER_WIDTH } from '@/shared/constants/layout';
import { appRoutes } from '@/shared/constants/routes';

const navigationItems = [
  {
    icon: DashboardRoundedIcon,
    label: 'Dashboard',
    to: appRoutes.dashboard
  },
  {
    icon: AddHomeWorkRoundedIcon,
    label: 'New Room',
    to: appRoutes.roomsNew
  },
  {
    icon: HistoryRoundedIcon,
    label: 'History',
    to: appRoutes.history
  },
  {
    icon: PersonRoundedIcon,
    label: 'Profile',
    to: appRoutes.profile
  },
  {
    icon: SettingsRoundedIcon,
    label: 'Settings',
    to: appRoutes.settings
  }
] as const;

export interface AppSidebarProps {
  readonly isDesktop: boolean;
}

export const AppSidebar = ({ isDesktop }: AppSidebarProps) => {
  const dispatch = useAppDispatch();
  const location = useLocation();
  const isOpen = useAppSelector(selectIsSidebarOpen);

  const drawerContent = (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <Stack spacing={0.5} sx={{ p: 3 }}>
        <Typography variant="overline">Workspace</Typography>
        <Typography variant="h5">Estimate Room</Typography>
        <Typography color="text.secondary" variant="body2">
          Page-by-page frontend scaffold with clean module boundaries.
        </Typography>
      </Stack>
      <List sx={{ px: 1.5, py: 2 }}>
        {navigationItems.map(({ icon: Icon, label, to }) => (
          <ListItemButton
            key={to}
            component={RouterLink}
            onClick={() => {
              if (!isDesktop) {
                dispatch(closeSidebar());
              }
            }}
            selected={location.pathname === to}
            sx={{
              position: 'relative',
              '&.Mui-selected::before': {
                backgroundColor: 'primary.main',
                borderRadius: 999,
                content: '""',
                height: 22,
                left: 10,
                position: 'absolute',
                top: '50%',
                transform: 'translateY(-50%)',
                width: 4
              }
            }}
            to={to}
          >
            <ListItemIcon sx={{ color: 'inherit', minWidth: 40 }}>
              <Icon />
            </ListItemIcon>
            <ListItemText primary={label} />
          </ListItemButton>
        ))}
      </List>
      <Box sx={{ mt: 'auto', p: 3 }}>
        <Typography color="text.secondary" variant="caption">
          API and WebSocket services are scaffolded in the shared layer for Go backend
          integration.
        </Typography>
      </Box>
    </Box>
  );

  return (
    <Drawer
      ModalProps={{ keepMounted: true }}
      onClose={() => dispatch(closeSidebar())}
      open={isDesktop ? true : isOpen}
      sx={{
        width: APP_DRAWER_WIDTH,
        flexShrink: 0,
        '& .MuiDrawer-paper': {
          boxSizing: 'border-box',
          width: APP_DRAWER_WIDTH
        }
      }}
      variant={isDesktop ? 'permanent' : 'temporary'}
    >
      {drawerContent}
    </Drawer>
  );
};
