import LogoutRoundedIcon from '@mui/icons-material/LogoutRounded';
import PersonRoundedIcon from '@mui/icons-material/PersonRounded';
import SettingsRoundedIcon from '@mui/icons-material/SettingsRounded';
import { ListItemIcon, Menu, MenuItem, Stack, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';

import {
  dashboardUserMenuPaperSx,
  dashboardUserMenuSummarySx
} from './DashboardUserMenu.styles';

export interface DashboardUserMenuProps {
  readonly anchorEl: HTMLElement | null;
  readonly displayName: string;
  readonly email: string | null;
  readonly isLoggingOut: boolean;
  readonly onClose: () => void;
  readonly onLogout: () => void;
}

export const DashboardUserMenu = ({
  anchorEl,
  displayName,
  email,
  isLoggingOut,
  onClose,
  onLogout
}: DashboardUserMenuProps) => (
  <Menu
    anchorEl={anchorEl}
    anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
    onClose={onClose}
    open={Boolean(anchorEl)}
    PaperProps={{
      elevation: 0,
      sx: dashboardUserMenuPaperSx
    }}
    transformOrigin={{ horizontal: 'right', vertical: 'top' }}
  >
    <Stack spacing={0.25} sx={dashboardUserMenuSummarySx}>
      <Typography variant="subtitle2">{displayName}</Typography>
      <Typography color="text.secondary" variant="caption">
        {email ?? 'No active session'}
      </Typography>
    </Stack>
    <MenuItem component={RouterLink} onClick={onClose} to={appRoutes.profile}>
      <ListItemIcon>
        <PersonRoundedIcon fontSize="small" />
      </ListItemIcon>
      Profile
    </MenuItem>
    <MenuItem component={RouterLink} onClick={onClose} to={appRoutes.settings}>
      <ListItemIcon>
        <SettingsRoundedIcon fontSize="small" />
      </ListItemIcon>
      Settings
    </MenuItem>
    <MenuItem
      disabled={isLoggingOut}
      onClick={() => {
        onClose();
        onLogout();
      }}
    >
      <ListItemIcon>
        <LogoutRoundedIcon fontSize="small" />
      </ListItemIcon>
      {isLoggingOut ? 'Logging out...' : 'Log out'}
    </MenuItem>
  </Menu>
);
