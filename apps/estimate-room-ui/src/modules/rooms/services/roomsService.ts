import type { Room } from '@/shared/types';

import type { NewRoomFormValues } from '../types';

const buildRoomModel = (input: NewRoomFormValues): Room => {
  const now = new Date().toISOString();

  return {
    createdAt: now,
    dimensions: {
      height: input.height,
      length: input.length,
      width: input.width
    },
    estimateStatus: 'draft',
    id: crypto.randomUUID(),
    name: input.name,
    teamId: input.teamId || null,
    updatedAt: now
  };
};

export const roomsService = {
  createRoom: async (input: NewRoomFormValues) => Promise.resolve(buildRoomModel(input)),
  getRoomPreview: (roomId: string): Room | null =>
    roomId
      ? {
          createdAt: '2026-03-16T14:30:00.000Z',
          dimensions: {
            height: 2.7,
            length: 6.8,
            width: 4.2
          },
          estimateStatus: 'queued',
          id: roomId,
          name: 'North Wing Suite',
          teamId: 'team-atlantic',
          updatedAt: '2026-03-18T09:15:00.000Z'
        }
      : null
};
