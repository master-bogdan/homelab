import ArrowForwardRoundedIcon from '@mui/icons-material/ArrowForwardRounded';
import MeetingRoomRoundedIcon from '@mui/icons-material/MeetingRoomRounded';
import NoteAltRoundedIcon from '@mui/icons-material/NoteAltRounded';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import {
  AppBox,
  AppButton,
  AppChip,
  AppPageState,
  AppStack,
  AppSurface,
  AppTypography
} from '@/shared/ui';

import type { DashboardSession } from '../../types';
import { formatRelativeTime, formatStatusLabel } from '../../utils';

import {
  recentRoomsCardActionLinkSx,
  recentRoomsCardArrowSx,
  recentRoomsCardItemIconSx,
  recentRoomsCardItemLinkSx,
  recentRoomsCardRootSx
} from './styles';

const getRoomRoute = (room: DashboardSession) =>
  room.status === 'ACTIVE'
    ? appRoutes.roomDetailsPath(room.id)
    : appRoutes.historyRoomPath(room.id);

export interface RecentRoomsCardProps {
  readonly onCreateRoom: () => void;
  readonly rooms: DashboardSession[];
}

export const RecentRoomsCard = ({ onCreateRoom, rooms }: RecentRoomsCardProps) => (
  <AppStack spacing={1.5}>
    <AppStack alignItems="center" direction="row" justifyContent="space-between">
      <AppTypography color="text.secondary" variant="overline">
        Recent Rooms
      </AppTypography>
      <AppTypography
        color="primary.main"
        component={RouterLink}
        sx={recentRoomsCardActionLinkSx}
        to={appRoutes.history}
        variant="caption"
      >
        View all
      </AppTypography>
    </AppStack>
    <AppSurface sx={recentRoomsCardRootSx(rooms.length === 0)}>
      {rooms.length === 0 ? (
        <AppPageState
          action={
            <AppButton onClick={onCreateRoom} variant="contained">
              Create room
            </AppButton>
          }
          description="New rooms will appear here after you start your first session."
          title="No recent rooms"
          visual={<MeetingRoomRoundedIcon color="disabled" fontSize="large" />}
        />
      ) : (
        <AppStack spacing={0.5}>
          {rooms.map((room) => (
            <AppBox
              component={RouterLink}
              key={`${room.id}-${room.lastActivityAt}`}
              sx={recentRoomsCardItemLinkSx}
              to={getRoomRoute(room)}
            >
              <AppStack alignItems="center" direction="row" spacing={1.5}>
                <AppSurface sx={recentRoomsCardItemIconSx(room.status === 'ACTIVE')}>
                  {room.status === 'ACTIVE' ? (
                    <MeetingRoomRoundedIcon fontSize="small" />
                  ) : (
                    <NoteAltRoundedIcon fontSize="small" />
                  )}
                </AppSurface>
                <AppStack spacing={0.25}>
                  <AppTypography variant="subtitle2">{room.name}</AppTypography>
                  <AppTypography color="text.secondary" variant="caption">
                    Modified {formatRelativeTime(room.lastActivityAt)}
                  </AppTypography>
                </AppStack>
              </AppStack>
              <AppStack alignItems="flex-end" spacing={0.75}>
                <AppChip
                  color={room.status === 'ACTIVE' ? 'success' : 'default'}
                  label={formatStatusLabel(room.status)}
                  size="small"
                  variant={room.status === 'ACTIVE' ? 'filled' : 'outlined'}
                />
                <AppTypography color="text.secondary" variant="caption">
                  {room.participantsCount} participants
                </AppTypography>
                <AppTypography color="primary.main" variant="caption">
                  Open
                  <ArrowForwardRoundedIcon fontSize="inherit" sx={recentRoomsCardArrowSx} />
                </AppTypography>
              </AppStack>
            </AppBox>
          ))}
        </AppStack>
      )}
    </AppSurface>
  </AppStack>
);
