import type { Room } from '@/modules/rooms/types';

export const selectRoomById = (rooms: Room[], roomId: string) =>
  rooms.find((room) => room.id === roomId) ?? null;
