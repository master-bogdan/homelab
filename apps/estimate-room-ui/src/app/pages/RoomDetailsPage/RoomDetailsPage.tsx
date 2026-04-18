import { useRoomDetailsPage } from '@/modules/rooms/hooks/useRoomDetailsPage';
import { formatDimensions, mapRoomStatusLabel } from '@/modules/rooms/utils';
import { AppChip, AppPageState, AppStack, AppTypography, SectionCard } from '@/shared/components';
import { formatDateTime } from '@/shared/utils';

export const RoomDetailsPage = () => {
  const { room, roomId } = useRoomDetailsPage();

  return (
    <SectionCard
      description="Placeholder detail view prepared for backend room data and estimate progress."
      title={`Room ${roomId}`}
    >
      {room ? (
        <AppStack spacing={2}>
          <AppTypography variant="h4">{room.name}</AppTypography>
          <AppStack direction="row" flexWrap="wrap" gap={1}>
            <AppChip color="primary" label={mapRoomStatusLabel(room.estimateStatus)} />
            <AppChip label={`Team: ${room.teamId ?? 'Unassigned'}`} variant="outlined" />
          </AppStack>
          <AppTypography color="text.secondary" variant="body2">
            Dimensions: {formatDimensions(room.dimensions)}
          </AppTypography>
          <AppTypography color="text.secondary" variant="body2">
            Updated: {formatDateTime(room.updatedAt)}
          </AppTypography>
        </AppStack>
      ) : (
        <AppPageState
          description="Connect the room details endpoint to populate dimensions, status, and team context."
          title="No room data yet"
          titleVariant="body1"
        />
      )}
    </SectionCard>
  );
};
