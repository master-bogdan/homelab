import type { Room } from '@/shared/types';

export interface NewRoomFormValues {
  readonly height: number;
  readonly length: number;
  readonly name: string;
  readonly teamId: string;
  readonly width: number;
}

export interface RoomPageData {
  readonly room: Room | null;
  readonly roomId: string;
}
