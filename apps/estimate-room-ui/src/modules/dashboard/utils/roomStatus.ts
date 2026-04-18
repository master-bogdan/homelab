import { DashboardRoomStatuses } from '../constants';
import type { DashboardRoomStatus } from '../types/status';

export const isActiveDashboardRoomStatus = (status: DashboardRoomStatus) =>
  status === DashboardRoomStatuses.ACTIVE;
