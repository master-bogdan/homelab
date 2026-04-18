import LogoutRoundedIcon from '@mui/icons-material/LogoutRounded';
import PersonRoundedIcon from '@mui/icons-material/PersonRounded';
import SettingsRoundedIcon from '@mui/icons-material/SettingsRounded';
import { Link as RouterLink } from 'react-router-dom';

import { AppRoutes } from '@/app/router/routePaths';
import {
  AppListItemIcon,
  AppMenu,
  AppMenuItem,
  AppStack,
  AppTypography
} from '@/shared/components';

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
  <AppMenu
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
    <AppStack spacing={0.25} sx={dashboardUserMenuSummarySx}>
      <AppTypography variant="subtitle2">{displayName}</AppTypography>
      <AppTypography color="text.secondary" variant="caption">
        {email ?? 'No active session'}
      </AppTypography>
    </AppStack>
    <AppMenuItem component={RouterLink} onClick={onClose} to={AppRoutes.PROFILE}>
      <AppListItemIcon>
        <PersonRoundedIcon fontSize="small" />
      </AppListItemIcon>
      Profile
    </AppMenuItem>
    <AppMenuItem component={RouterLink} onClick={onClose} to={AppRoutes.SETTINGS}>
      <AppListItemIcon>
        <SettingsRoundedIcon fontSize="small" />
      </AppListItemIcon>
      Settings
    </AppMenuItem>
    <AppMenuItem disabled={isLoggingOut} onClick={onLogout}>
      <AppListItemIcon>
        <LogoutRoundedIcon fontSize="small" />
      </AppListItemIcon>
      {isLoggingOut ? 'Logging out...' : 'Log out'}
    </AppMenuItem>
  </AppMenu>
);
