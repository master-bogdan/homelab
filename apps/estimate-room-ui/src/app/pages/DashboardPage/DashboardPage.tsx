import WarningAmberRoundedIcon from '@mui/icons-material/WarningAmberRounded';
import { useNavigate } from 'react-router-dom';

import { AppRoutes } from '@/app/router/routePaths';
import { AppAlert, AppBox, AppButton, AppPageState, AppStack, AppSurface } from '@/shared/components';

import {
  DashboardHeroCard,
  LedgerCard,
  RecentRoomsCard,
  TeamsCard
} from '@/modules/dashboard/components';
import { DashboardLoadStatuses } from '@/modules/dashboard';
import { useDashboardActions, useDashboardPage } from '@/modules/dashboard/hooks';
import {
  dashboardPageActiveGridSx,
  dashboardPageNoActiveSx,
  dashboardPageSectionsGridSx,
  dashboardPageStateCardSx
} from './DashboardPage.styles';

export const DashboardPage = () => {
  const navigate = useNavigate();
  const { openCreateRoom } = useDashboardActions();
  const { data, errorMessage, retry, status } = useDashboardPage();

  if (status === DashboardLoadStatuses.LOADING) {
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

  if (status === DashboardLoadStatuses.ERROR || !data) {
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
            onOpenRoom={(roomId) => navigate(AppRoutes.ROOM_DETAILS_PATH(roomId))}
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
              onOpenRoom={(roomId) => navigate(AppRoutes.ROOM_DETAILS_PATH(roomId))}
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
