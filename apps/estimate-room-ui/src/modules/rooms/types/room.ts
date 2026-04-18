export interface RoomDimensions {
  readonly height: number;
  readonly length: number;
  readonly width: number;
}

export interface Room {
  readonly createdAt: string;
  readonly dimensions: RoomDimensions;
  readonly estimateStatus: 'draft' | 'queued' | 'completed';
  readonly id: string;
  readonly name: string;
  readonly teamId: string | null;
  readonly updatedAt: string;
}
