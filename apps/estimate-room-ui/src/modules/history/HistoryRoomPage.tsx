import { Chip, Stack, Typography } from '@mui/material';

import { AppPageState, SectionCard } from '@/shared/ui';
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
      <Stack spacing={2}>
        {entries.length ? (
          entries.map((entry) => (
            <Stack key={entry.id} spacing={1}>
              <Stack alignItems="center" direction="row" spacing={1}>
                <Chip color={mapHistoryStatusColor(entry.status)} label={entry.status} />
                <Typography variant="body2">{formatDateTime(entry.capturedAt)}</Typography>
              </Stack>
              <Typography color="text.secondary" variant="body2">
                Submitted by {entry.submittedBy}
              </Typography>
            </Stack>
          ))
        ) : (
          <AppPageState
            description="When room activity is connected, event history and processing status will appear here."
            title="No history entries yet"
            titleVariant="body1"
          />
        )}
      </Stack>
    </SectionCard>
  );
};
