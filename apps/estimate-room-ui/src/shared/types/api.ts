export interface ApiErrorItem {
  readonly detail: string;
  readonly pointer: string;
}

export interface ApiError {
  readonly detail?: string;
  readonly errors?: ApiErrorItem[];
  readonly instance?: string;
  readonly code?: string;
  readonly details?: Record<string, unknown>;
  readonly message: string;
  readonly status: number;
  readonly title?: string;
  readonly type?: string;
}

export interface PaginatedResponse<TItem> {
  readonly items: TItem[];
  readonly page: number;
  readonly pageSize: number;
  readonly total: number;
}
