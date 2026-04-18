import { Link as RouterLink } from 'react-router-dom';

import { AppRoutes } from '@/app/router/routePaths';
import { useHistoryPage } from '@/modules/history/hooks/useHistoryPage';
import { mapHistoryStatusColor } from '@/modules/history/utils';
import { AppChip, AppStack, AppTypography, SectionCard } from '@/shared/components';
import { formatDateTime } from '@/shared/utils';

export const HistoryPage = () => {
  const { entries } = useHistoryPage();

  return (
    <SectionCard
      description="Estimate submissions and processing checkpoints will land here once the backend history endpoints are connected."
      title="History"
    >
      <AppStack spacing={2}>
        {entries.map((entry) => (
          <SectionCard
            key={entry.id}
            action={
              <AppTypography
                color="primary"
                component={RouterLink}
                to={AppRoutes.HISTORY_ROOM_PATH(entry.roomId)}
                variant="body2"
              >
                View room history
              </AppTypography>
            }
            description={formatDateTime(entry.capturedAt)}
            title={`Room ${entry.roomId}`}
          >
            <AppStack alignItems="center" direction="row" spacing={1}>
              <AppChip color={mapHistoryStatusColor(entry.status)} label={entry.status} />
              <AppTypography color="text.secondary" variant="body2">
                Submitted by {entry.submittedBy}
              </AppTypography>
            </AppStack>
          </SectionCard>
        ))}
      </AppStack>
    </SectionCard>
  );
};
