export interface RoomEstimateHistoryEntry {
  readonly capturedAt: string;
  readonly id: string;
  readonly roomId: string;
  readonly status: 'queued' | 'processed' | 'failed';
  readonly submittedBy: string;
}
