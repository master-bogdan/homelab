import type { SvgIconComponent } from '@mui/icons-material';
import DashboardRoundedIcon from '@mui/icons-material/DashboardRounded';
import HistoryRoundedIcon from '@mui/icons-material/HistoryRounded';
import PersonRoundedIcon from '@mui/icons-material/PersonRounded';
import SettingsRoundedIcon from '@mui/icons-material/SettingsRounded';

import { AppRoutes } from '@/shared/constants/routes';

export interface DashboardLayoutNavigationItem {
  readonly icon: SvgIconComponent;
  readonly label: string;
  readonly to: string;
}

export const primaryNavigationItems = [
  {
    icon: DashboardRoundedIcon,
    label: 'Dashboard',
    to: AppRoutes.DASHBOARD
  },
  {
    icon: HistoryRoundedIcon,
    label: 'History',
    to: AppRoutes.HISTORY
  }
] as const satisfies readonly DashboardLayoutNavigationItem[];

export const secondaryNavigationItems = [
  {
    icon: PersonRoundedIcon,
    label: 'Profile',
    to: AppRoutes.PROFILE
  },
  {
    icon: SettingsRoundedIcon,
    label: 'Settings',
    to: AppRoutes.SETTINGS
  }
] as const satisfies readonly DashboardLayoutNavigationItem[];
