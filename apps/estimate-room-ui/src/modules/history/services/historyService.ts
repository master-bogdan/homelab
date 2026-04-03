import type { RoomEstimateHistoryEntry } from '@/shared/types';

const historyEntries: RoomEstimateHistoryEntry[] = [
  {
    capturedAt: '2026-03-18T15:10:00.000Z',
    id: 'hist-001',
    roomId: 'room-101',
    status: 'processed',
    submittedBy: 'Alex Morgan'
  },
  {
    capturedAt: '2026-03-18T14:05:00.000Z',
    id: 'hist-002',
    roomId: 'room-102',
    status: 'queued',
    submittedBy: 'Priya Sharma'
  },
  {
    capturedAt: '2026-03-17T18:42:00.000Z',
    id: 'hist-003',
    roomId: 'room-103',
    status: 'failed',
    submittedBy: 'Luis Rivera'
  }
];

export const historyService = {
  getHistoryEntries: () => historyEntries,
  getHistoryForRoom: (roomId: string) =>
    historyEntries.filter((entry) => entry.roomId === roomId)
};
