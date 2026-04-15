import { AppChip, AppPageState, AppStack, AppTypography, SectionCard } from '@/shared/ui';
import { formatDateTime } from '@/shared/utils';

import { useHistoryRoomPage } from './hooks/useHistoryRoomPage';
import { mapHistoryStatusColor } from './utils';

export const HistoryRoomPage = () => {
  const { entries, roomId } = useHistoryRoomPage();

  return (
    <SectionCard
      description="Backend event history for a single room will expand here with job progress, retries, and operator notes."
      title={`Room History ${roomId}`}
    >
      <AppStack spacing={2}>
        {entries.length ? (
          entries.map((entry) => (
            <AppStack key={entry.id} spacing={1}>
              <AppStack alignItems="center" direction="row" spacing={1}>
                <AppChip color={mapHistoryStatusColor(entry.status)} label={entry.status} />
                <AppTypography variant="body2">{formatDateTime(entry.capturedAt)}</AppTypography>
              </AppStack>
              <AppTypography color="text.secondary" variant="body2">
                Submitted by {entry.submittedBy}
              </AppTypography>
            </AppStack>
          ))
        ) : (
          <AppPageState
            description="When room activity is connected, event history and processing status will appear here."
            title="No history entries yet"
            titleVariant="body1"
          />
        )}
      </AppStack>
    </SectionCard>
  );
};
