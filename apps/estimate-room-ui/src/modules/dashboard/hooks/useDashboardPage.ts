import { usePageTitle } from '@/shared/hooks';

import { dashboardService } from '../services/dashboardService';

export const useDashboardPage = () => {
  usePageTitle('Dashboard');

  return {
    metrics: dashboardService.getOverviewMetrics(),
    workstreams: dashboardService.getWorkstreams()
  };
};
