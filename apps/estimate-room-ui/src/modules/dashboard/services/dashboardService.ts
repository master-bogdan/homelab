import type { DashboardMetric, DashboardWorkstream } from '../types';

const overviewMetrics: DashboardMetric[] = [
  {
    helperText: 'Rooms waiting for backend scoring or post-processing.',
    id: 'queued-estimates',
    label: 'Queued Estimates',
    value: '12'
  },
  {
    helperText: 'Projects touched within the last 24 hours.',
    id: 'active-rooms',
    label: 'Active Rooms',
    value: '28'
  },
  {
    helperText: 'Current spread across field and operations teams.',
    id: 'connected-teams',
    label: 'Connected Teams',
    value: '7'
  }
];

const workstreams: DashboardWorkstream[] = [
  {
    description: 'Finalize DTO mapping for room creation and estimate requests.',
    id: 'backend-contract',
    title: 'Backend Contract Review'
  },
  {
    description: 'Prepare a live stream consumer for long-running estimate jobs.',
    id: 'ws-jobs',
    title: 'WebSocket Progress Updates'
  },
  {
    description: 'Sequence screen-by-screen implementation from generated designs.',
    id: 'design-rollout',
    title: 'Design Integration Queue'
  }
];

export const dashboardService = {
  getOverviewMetrics: () => overviewMetrics,
  getWorkstreams: () => workstreams
};
