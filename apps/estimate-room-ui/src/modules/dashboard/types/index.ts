export interface DashboardMetric {
  readonly helperText: string;
  readonly id: string;
  readonly label: string;
  readonly value: string;
}

export interface DashboardWorkstream {
  readonly id: string;
  readonly title: string;
  readonly description: string;
}
