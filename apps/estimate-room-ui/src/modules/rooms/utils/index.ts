import type { Room, RoomDimensions } from '@/modules/rooms/types';

export const formatDimensions = (dimensions: RoomDimensions) =>
  `${dimensions.length}m × ${dimensions.width}m × ${dimensions.height}m`;

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
