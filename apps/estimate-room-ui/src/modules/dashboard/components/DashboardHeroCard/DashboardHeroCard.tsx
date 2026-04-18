import AddCircleRoundedIcon from '@mui/icons-material/AddCircleRounded';
import ArrowForwardRoundedIcon from '@mui/icons-material/ArrowForwardRounded';
import MeetingRoomRoundedIcon from '@mui/icons-material/MeetingRoomRounded';
import RadioButtonCheckedRoundedIcon from '@mui/icons-material/RadioButtonCheckedRounded';

import {
  AppAvatar,
  AppBox,
  AppButton,
  AppProgress,
  AppStack,
  AppSurface,
  AppTypography
} from '@/shared/components';

import type { DashboardActiveRoom } from '../../types';
import { formatRelativeTime, formatStatusLabel, getInitials } from '../../utils';

import {
  dashboardHeroCardEmptyAvatarSx,
  dashboardHeroCardEmptyRootSx,
  dashboardHeroCardMemberOverflowSx,
  dashboardHeroCardMetricCardSx,
  dashboardHeroCardMetricGridSx,
  dashboardHeroCardParticipantAvatarSx,
  dashboardHeroCardProgressBarSx,
  dashboardHeroCardRootSx
} from './styles';

export interface DashboardHeroCardProps {
  readonly onCreateRoom: () => void;
  readonly onOpenRoom: (roomId: string) => void;
  readonly room: DashboardActiveRoom | null;
}

export const DashboardHeroCard = ({
  onCreateRoom,
  onOpenRoom,
  room
}: DashboardHeroCardProps) => {
  if (!room) {
    return (
      <AppSurface sx={dashboardHeroCardEmptyRootSx}>
        <AppStack alignItems="center" spacing={2.5}>
          <AppAvatar sx={dashboardHeroCardEmptyAvatarSx}>
            <MeetingRoomRoundedIcon />
          </AppAvatar>
          <AppStack spacing={1.5}>
            <AppTypography variant="h4">No active rooms yet</AppTypography>
            <AppTypography color="text.secondary" maxWidth={320} variant="body2">
              Start your first collaborative architectural refinement session with your team.
            </AppTypography>
          </AppStack>
          <AppButton
            onClick={onCreateRoom}
            startIcon={<AddCircleRoundedIcon />}
            variant="contained"
          >
            Create your first room
          </AppButton>
        </AppStack>
      </AppSurface>
    );
  }

  const visibleParticipants = room.participants.slice(0, 4);
  const remainingParticipants = room.participants.length - visibleParticipants.length;
  const taskProgress =
    room.tasksCount === 0
      ? 0
      : Math.round((room.estimatedTasksCount / room.tasksCount) * 100);
  const statusLabel = formatStatusLabel(room.currentTaskStatus ?? room.status);

  return (
    <AppSurface sx={dashboardHeroCardRootSx}>
      <AppStack
        alignItems={{ sm: 'flex-start' }}
        direction={{ sm: 'row' }}
        justifyContent="space-between"
        spacing={2}
      >
        <AppStack spacing={1.5}>
          <AppTypography color="primary.main" variant="overline">
            Active Session
          </AppTypography>
          <AppStack spacing={0.5}>
            <AppTypography variant="h4">{room.name}</AppTypography>
            <AppTypography color="text.secondary" variant="body2">
              Room code {room.code} • Active {formatRelativeTime(room.lastActivityAt)}
            </AppTypography>
          </AppStack>
        </AppStack>
        <AppStack direction="row" sx={{ ml: { sm: 'auto' } }}>
          {visibleParticipants.map((participant, index) => (
            <AppAvatar
              key={participant.id}
              src={participant.avatarUrl ?? undefined}
              sx={dashboardHeroCardParticipantAvatarSx(index)}
            >
              {getInitials(participant.displayName)}
            </AppAvatar>
          ))}
          {remainingParticipants > 0 ? (
            <AppAvatar sx={dashboardHeroCardMemberOverflowSx}>+{remainingParticipants}</AppAvatar>
          ) : null}
        </AppStack>
      </AppStack>
      <AppBox sx={dashboardHeroCardMetricGridSx}>
        <AppSurface sx={dashboardHeroCardMetricCardSx}>
          <AppTypography color="text.secondary" variant="overline">
            Current Task
          </AppTypography>
          <AppTypography sx={{ mt: 0.5 }} variant="subtitle2">
            {room.currentTaskTitle ?? 'No task selected'}
          </AppTypography>
        </AppSurface>
        <AppSurface sx={dashboardHeroCardMetricCardSx}>
          <AppTypography color="text.secondary" variant="overline">
            Participants
          </AppTypography>
          <AppTypography sx={{ mt: 0.5 }} variant="subtitle2">
            {room.participants.length} online
          </AppTypography>
        </AppSurface>
        <AppSurface sx={dashboardHeroCardMetricCardSx}>
          <AppTypography color="text.secondary" variant="overline">
            Status
          </AppTypography>
          <AppStack alignItems="center" direction="row" spacing={1} sx={{ mt: 0.5 }}>
            <RadioButtonCheckedRoundedIcon color="success" fontSize="small" />
            <AppTypography variant="subtitle2">{statusLabel}</AppTypography>
          </AppStack>
        </AppSurface>
      </AppBox>
      <AppStack spacing={1.25}>
        <AppStack alignItems="center" direction="row" justifyContent="space-between">
          <AppTypography color="text.secondary" variant="overline">
            Estimation Progress
          </AppTypography>
          <AppTypography color="primary.main" variant="caption">
            {room.estimatedTasksCount} / {room.tasksCount || 0} tasks estimated
          </AppTypography>
        </AppStack>
        <AppProgress kind="linear" sx={dashboardHeroCardProgressBarSx} value={taskProgress} variant="determinate" />
      </AppStack>
      <AppButton
        endIcon={<ArrowForwardRoundedIcon />}
        onClick={() => onOpenRoom(room.id)}
        sx={{ mt: 'auto' }}
        variant="contained"
      >
        Enter Room
      </AppButton>
    </AppSurface>
  );
};
