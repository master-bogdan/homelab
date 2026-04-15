import WarningAmberRoundedIcon from '@mui/icons-material/WarningAmberRounded';
import { Alert, Box, Paper, Stack } from '@mui/material';
import { useNavigate } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppButton, AppPageState } from '@/shared/ui';

import {
  DashboardHeroCard,
  LedgerCard,
  RecentRoomsCard,
  TeamsCard
} from './components';
import {
  dashboardPageActiveGridSx,
  dashboardPageNoActiveSx,
  dashboardPageSectionsGridSx,
  dashboardPageStateCardSx
} from './DashboardPage.styles';
import { useDashboardActions } from './hooks/useDashboardActions';
import { useDashboardPage } from './hooks/useDashboardPage';

export const DashboardPage = () => {
  const navigate = useNavigate();
  const { openCreateRoom } = useDashboardActions();
  const { data, errorMessage, retry, status } = useDashboardPage();

  if (status === 'loading') {
    return (
      <Paper elevation={0} sx={dashboardPageStateCardSx}>
        <AppPageState
          description="Pulling your session history, teams, and architect ledger."
          isLoading
          title="Loading dashboard"
        />
      </Paper>
    );
  }

  if (status === 'error' || !data) {
    return (
      <Paper elevation={0} sx={dashboardPageStateCardSx}>
        <AppPageState
          action={
            <AppButton onClick={retry} variant="contained">
              Retry
            </AppButton>
          }
          description={errorMessage}
          title="Dashboard unavailable"
          visual={<WarningAmberRoundedIcon color="warning" fontSize="large" />}
        />
      </Paper>
    );
  }

  return (
    <Stack spacing={4}>
      {data.activeRoomError && data.view !== 'active' ? (
        <Alert severity="warning">{data.activeRoomError.message}</Alert>
      ) : null}
      {data.view === 'active' ? (
        <Box sx={dashboardPageActiveGridSx}>
          <DashboardHeroCard
            onCreateRoom={openCreateRoom}
            onOpenRoom={(roomId) => navigate(appRoutes.roomDetailsPath(roomId))}
            room={data.activeRoom}
          />
          <RecentRoomsCard
            onCreateRoom={openCreateRoom}
            rooms={data.recentRooms.filter((room) => room.id !== data.activeRoom?.id)}
          />
        </Box>
      ) : (
        <Box sx={dashboardPageNoActiveSx}>
          {data.view === 'noActive' ? (
            <RecentRoomsCard onCreateRoom={openCreateRoom} rooms={data.recentRooms} />
          ) : (
            <DashboardHeroCard
              onCreateRoom={openCreateRoom}
              onOpenRoom={(roomId) => navigate(appRoutes.roomDetailsPath(roomId))}
              room={null}
            />
          )}
        </Box>
      )}
      <Box sx={dashboardPageSectionsGridSx}>
        <TeamsCard
          errorMessage={data.teamsError?.message ?? null}
          onRetry={retry}
          teams={data.teams}
        />
        <LedgerCard
          errorMessage={data.ledgerError?.message ?? null}
          ledger={data.ledger}
          onRetry={retry}
        />
      </Box>
    </Stack>
  );
};
