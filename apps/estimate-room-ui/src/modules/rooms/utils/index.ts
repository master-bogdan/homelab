import type { Room } from '@/shared/types';

export const mapRoomStatusLabel = (status: Room['estimateStatus']) => {
  switch (status) {
    case 'completed':
      return 'Completed';
    case 'queued':
      return 'Queued';
    default:
      return 'Draft';
  }
};
