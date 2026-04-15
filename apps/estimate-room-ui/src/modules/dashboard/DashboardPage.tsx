import WarningAmberRoundedIcon from '@mui/icons-material/WarningAmberRounded';
import { useNavigate } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppAlert, AppBox, AppButton, AppPageState, AppStack, AppSurface } from '@/shared/ui';

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
      <AppSurface sx={dashboardPageStateCardSx}>
        <AppPageState
          description="Pulling your session history, teams, and architect ledger."
          isLoading
          title="Loading dashboard"
        />
      </AppSurface>
    );
  }

  if (status === 'error' || !data) {
    return (
      <AppSurface sx={dashboardPageStateCardSx}>
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
      </AppSurface>
    );
  }

  return (
    <AppStack spacing={4}>
      {data.activeRoomError && data.view !== 'active' ? (
        <AppAlert severity="warning">{data.activeRoomError.message}</AppAlert>
      ) : null}
      {data.view === 'active' ? (
        <AppBox sx={dashboardPageActiveGridSx}>
          <DashboardHeroCard
            onCreateRoom={openCreateRoom}
            onOpenRoom={(roomId) => navigate(appRoutes.roomDetailsPath(roomId))}
            room={data.activeRoom}
          />
          <RecentRoomsCard
            onCreateRoom={openCreateRoom}
            rooms={data.recentRooms.filter((room) => room.id !== data.activeRoom?.id)}
          />
        </AppBox>
      ) : (
        <AppBox sx={dashboardPageNoActiveSx}>
          {data.view === 'noActive' ? (
            <RecentRoomsCard onCreateRoom={openCreateRoom} rooms={data.recentRooms} />
          ) : (
            <DashboardHeroCard
              onCreateRoom={openCreateRoom}
              onOpenRoom={(roomId) => navigate(appRoutes.roomDetailsPath(roomId))}
              room={null}
            />
          )}
        </AppBox>
      )}
      <AppBox sx={dashboardPageSectionsGridSx}>
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
      </AppBox>
    </AppStack>
  );
};
