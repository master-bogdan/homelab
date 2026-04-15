import type { SvgIconComponent } from '@mui/icons-material';
import DashboardRoundedIcon from '@mui/icons-material/DashboardRounded';
import HistoryRoundedIcon from '@mui/icons-material/HistoryRounded';
import PersonRoundedIcon from '@mui/icons-material/PersonRounded';
import SettingsRoundedIcon from '@mui/icons-material/SettingsRounded';

import { appRoutes } from '@/shared/constants/routes';

export interface DashboardLayoutNavigationItem {
  readonly icon: SvgIconComponent;
  readonly label: string;
  readonly to: string;
}

export const primaryNavigationItems = [
  {
    icon: DashboardRoundedIcon,
    label: 'Dashboard',
    to: appRoutes.dashboard
  },
  {
    icon: HistoryRoundedIcon,
    label: 'History',
    to: appRoutes.history
  }
] as const satisfies readonly DashboardLayoutNavigationItem[];

export const secondaryNavigationItems = [
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
] as const satisfies readonly DashboardLayoutNavigationItem[];
