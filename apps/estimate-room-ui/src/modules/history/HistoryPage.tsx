import { Chip, Stack, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { SectionCard } from '@/shared/ui';
import { formatDateTime } from '@/shared/utils';

import { useHistoryPage } from './hooks/useHistoryPage';
import { mapHistoryStatusColor } from './utils';

export const HistoryPage = () => {
  const { entries } = useHistoryPage();

  return (
    <SectionCard
      description="Estimate submissions and processing checkpoints will land here once the backend history endpoints are connected."
      title="History"
    >
      <Stack spacing={2}>
        {entries.map((entry) => (
          <SectionCard
            key={entry.id}
            action={
              <Typography
                color="primary"
                component={RouterLink}
                to={appRoutes.historyRoomPath(entry.roomId)}
                variant="body2"
              >
                View room history
              </Typography>
            }
            description={formatDateTime(entry.capturedAt)}
            title={`Room ${entry.roomId}`}
          >
            <Stack alignItems="center" direction="row" spacing={1}>
              <Chip color={mapHistoryStatusColor(entry.status)} label={entry.status} />
              <Typography color="text.secondary" variant="body2">
                Submitted by {entry.submittedBy}
              </Typography>
            </Stack>
          </SectionCard>
        ))}
      </Stack>
    </SectionCard>
  );
};
