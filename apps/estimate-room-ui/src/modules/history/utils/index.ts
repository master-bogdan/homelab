import type { RoomEstimateHistoryEntry } from '@/modules/history/types';

export const mapHistoryStatusColor = (status: RoomEstimateHistoryEntry['status']) => {
  switch (status) {
    case 'processed':
      return 'success' as const;
    case 'failed':
      return 'error' as const;
    default:
      return 'warning' as const;
  }
};
