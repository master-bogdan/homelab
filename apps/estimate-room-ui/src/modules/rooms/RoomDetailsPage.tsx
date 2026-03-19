import { Chip, Stack, Typography } from '@mui/material';

import { SectionCard } from '@/shared/ui';
import { formatDateTime, formatDimensions } from '@/shared/utils';

import { useRoomDetailsPage } from './hooks/useRoomDetailsPage';
import { mapRoomStatusLabel } from './utils';

export const RoomDetailsPage = () => {
  const { room, roomId } = useRoomDetailsPage();

  return (
    <SectionCard
      description="Placeholder detail view prepared for backend room data and estimate progress."
      title={`Room ${roomId}`}
    >
      {room ? (
        <Stack spacing={2}>
          <Typography variant="h4">{room.name}</Typography>
          <Stack direction="row" flexWrap="wrap" gap={1}>
            <Chip color="primary" label={mapRoomStatusLabel(room.estimateStatus)} />
            <Chip label={`Team: ${room.teamId ?? 'Unassigned'}`} variant="outlined" />
          </Stack>
          <Typography color="text.secondary" variant="body2">
            Dimensions: {formatDimensions(room.dimensions)}
          </Typography>
          <Typography color="text.secondary" variant="body2">
            Updated: {formatDateTime(room.updatedAt)}
          </Typography>
        </Stack>
      ) : (
        <Typography color="text.secondary" variant="body2">
          No room data is available for this route yet.
        </Typography>
      )}
    </SectionCard>
  );
};
