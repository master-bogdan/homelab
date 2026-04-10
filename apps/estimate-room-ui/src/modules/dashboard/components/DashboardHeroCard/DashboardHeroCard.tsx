import AddCircleRoundedIcon from '@mui/icons-material/AddCircleRounded';
import ArrowForwardRoundedIcon from '@mui/icons-material/ArrowForwardRounded';
import MeetingRoomRoundedIcon from '@mui/icons-material/MeetingRoomRounded';
import RadioButtonCheckedRoundedIcon from '@mui/icons-material/RadioButtonCheckedRounded';
import {
  Avatar,
  Box,
  LinearProgress,
  Paper,
  Stack,
  Typography
} from '@mui/material';

import { AppButton } from '@/shared/ui';

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
      <Paper elevation={0} sx={dashboardHeroCardEmptyRootSx}>
        <Stack alignItems="center" spacing={2.5}>
          <Avatar sx={dashboardHeroCardEmptyAvatarSx}>
            <MeetingRoomRoundedIcon />
          </Avatar>
          <Stack spacing={1.5}>
            <Typography variant="h4">No active rooms yet</Typography>
            <Typography color="text.secondary" maxWidth={320} variant="body2">
              Start your first collaborative architectural refinement session with your team.
            </Typography>
          </Stack>
          <AppButton
            onClick={onCreateRoom}
            startIcon={<AddCircleRoundedIcon />}
            variant="contained"
          >
            Create your first room
          </AppButton>
        </Stack>
      </Paper>
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
    <Paper elevation={0} sx={dashboardHeroCardRootSx}>
      <Stack
        alignItems={{ sm: 'flex-start' }}
        direction={{ sm: 'row' }}
        justifyContent="space-between"
        spacing={2}
      >
        <Stack spacing={1.5}>
          <Typography color="primary.main" variant="overline">
            Active Session
          </Typography>
          <Stack spacing={0.5}>
            <Typography variant="h4">{room.name}</Typography>
            <Typography color="text.secondary" variant="body2">
              Room code {room.code} • Active {formatRelativeTime(room.lastActivityAt)}
            </Typography>
          </Stack>
        </Stack>
        <Stack direction="row" sx={{ ml: { sm: 'auto' } }}>
          {visibleParticipants.map((participant, index) => (
            <Avatar
              key={participant.id}
              src={participant.avatarUrl ?? undefined}
              sx={dashboardHeroCardParticipantAvatarSx(index)}
            >
              {getInitials(participant.displayName)}
            </Avatar>
          ))}
          {remainingParticipants > 0 ? (
            <Avatar sx={dashboardHeroCardMemberOverflowSx}>+{remainingParticipants}</Avatar>
          ) : null}
        </Stack>
      </Stack>
      <Box sx={dashboardHeroCardMetricGridSx}>
        <Paper elevation={0} sx={dashboardHeroCardMetricCardSx}>
          <Typography color="text.secondary" variant="overline">
            Current Task
          </Typography>
          <Typography sx={{ mt: 0.5 }} variant="subtitle2">
            {room.currentTaskTitle ?? 'No task selected'}
          </Typography>
        </Paper>
        <Paper elevation={0} sx={dashboardHeroCardMetricCardSx}>
          <Typography color="text.secondary" variant="overline">
            Participants
          </Typography>
          <Typography sx={{ mt: 0.5 }} variant="subtitle2">
            {room.participants.length} online
          </Typography>
        </Paper>
        <Paper elevation={0} sx={dashboardHeroCardMetricCardSx}>
          <Typography color="text.secondary" variant="overline">
            Status
          </Typography>
          <Stack alignItems="center" direction="row" spacing={1} sx={{ mt: 0.5 }}>
            <RadioButtonCheckedRoundedIcon color="success" fontSize="small" />
            <Typography variant="subtitle2">{statusLabel}</Typography>
          </Stack>
        </Paper>
      </Box>
      <Stack spacing={1.25}>
        <Stack alignItems="center" direction="row" justifyContent="space-between">
          <Typography color="text.secondary" variant="overline">
            Estimation Progress
          </Typography>
          <Typography color="primary.main" variant="caption">
            {room.estimatedTasksCount} / {room.tasksCount || 0} tasks estimated
          </Typography>
        </Stack>
        <LinearProgress sx={dashboardHeroCardProgressBarSx} value={taskProgress} variant="determinate" />
      </Stack>
      <AppButton
        endIcon={<ArrowForwardRoundedIcon />}
        onClick={() => onOpenRoom(room.id)}
        sx={{ mt: 'auto' }}
        variant="contained"
      >
        Enter Room
      </AppButton>
    </Paper>
  );
};
