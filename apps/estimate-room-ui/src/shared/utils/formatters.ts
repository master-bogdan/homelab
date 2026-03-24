import type { RoomDimensions } from '@/shared/types';

export const formatDimensions = (dimensions: RoomDimensions) =>
  `${dimensions.length}m × ${dimensions.width}m × ${dimensions.height}m`;

export const formatDateTime = (value: string) =>
  new Intl.DateTimeFormat('en-US', {
    dateStyle: 'medium',
    timeStyle: 'short'
  }).format(new Date(value));
