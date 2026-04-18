import type { DashboardCreateRoomFormValues } from '../types/forms';

export const DashboardCreateRoomLimits = {
  ROOM_NAME_MAX_LENGTH: 100
} as const;

export const DashboardCreateRoomDefaultValues: DashboardCreateRoomFormValues = {
  createShareLink: true,
  deckKey: 'fibonacci',
  inviteEmails: '',
  inviteTeamId: '',
  name: ''
};
