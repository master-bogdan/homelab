import ArrowForwardRoundedIcon from '@mui/icons-material/ArrowForwardRounded';
import { Box, Stack, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppButton, SectionCard } from '@/shared/ui';

import { useDashboardPage } from './hooks/useDashboardPage';

export const DashboardPage = () => {
  const { metrics, workstreams } = useDashboardPage();

  return (
    <Stack spacing={3}>
      <SectionCard
        action={
          <AppButton
            component={RouterLink}
            endIcon={<ArrowForwardRoundedIcon />}
            to={appRoutes.roomsNew}
            variant="contained"
          >
            Create room
          </AppButton>
        }
        description="Use this shell as the landing page for estimate operations, activity feeds, and backend status."
        title="Dashboard"
      >
        <Typography color="text.secondary" variant="body2">
          The initial scaffold keeps the module lightweight while leaving clear room for
          real API-backed widgets.
        </Typography>
      </SectionCard>

      <Box
        sx={{
          display: 'grid',
          gap: 3,
          gridTemplateColumns: {
            xs: 'repeat(1, minmax(0, 1fr))',
            md: 'repeat(2, minmax(0, 1fr))',
            xl: 'repeat(3, minmax(0, 1fr))'
          }
        }}
      >
        {metrics.map((metric) => (
          <Box key={metric.id}>
            <SectionCard description={metric.helperText} title={metric.label}>
              <Typography variant="h3">{metric.value}</Typography>
            </SectionCard>
          </Box>
        ))}
      </Box>

      <SectionCard
        description="Suggested next areas for page-by-page implementation."
        title="Delivery Focus"
      >
        <Stack spacing={2}>
          {workstreams.map((workstream) => (
            <Stack key={workstream.id} spacing={0.5}>
              <Typography variant="h6">{workstream.title}</Typography>
              <Typography color="text.secondary" variant="body2">
                {workstream.description}
              </Typography>
            </Stack>
          ))}
        </Stack>
      </SectionCard>
    </Stack>
  );
};
