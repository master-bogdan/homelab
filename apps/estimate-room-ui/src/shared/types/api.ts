export interface ApiError {
  readonly code?: string;
  readonly details?: Record<string, unknown>;
  readonly message: string;
  readonly status: number;
}

export interface PaginatedResponse<TItem> {
  readonly items: TItem[];
  readonly page: number;
  readonly pageSize: number;
  readonly total: number;
}
