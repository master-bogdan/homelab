import type { Room } from '@/shared/types';

export const selectRoomById = (rooms: Room[], roomId: string) =>
  rooms.find((room) => room.id === roomId) ?? null;
