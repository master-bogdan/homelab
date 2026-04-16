import GroupRoundedIcon from '@mui/icons-material/GroupRounded';
import { Link as RouterLink } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';
import {
  AppBox,
  AppButton,
  AppPageState,
  AppStack,
  AppSurface,
  AppTypography
} from '@/shared/ui';

import type { DashboardTeamSummary } from '../../types';
import { formatRelativeTime, getInitials } from '../../utils';

import {
  teamsCardAvatarSx,
  teamsCardItemLinkSx,
  teamsCardRootSx
} from './styles';

export interface TeamsCardProps {
  readonly errorMessage: string | null;
  readonly onRetry: () => void;
  readonly teams: DashboardTeamSummary[];
}

export const TeamsCard = ({ errorMessage, onRetry, teams }: TeamsCardProps) => (
  <AppStack spacing={1.5}>
    <AppTypography color="text.secondary" variant="overline">
      Your Teams
    </AppTypography>
    <AppSurface sx={teamsCardRootSx(teams.length === 0 || Boolean(errorMessage))}>
      {errorMessage ? (
        <AppPageState
          action={
            <AppButton onClick={onRetry} variant="contained">
              Retry
            </AppButton>
          }
          description={errorMessage}
          title="Teams are temporarily unavailable"
          visual={<GroupRoundedIcon color="disabled" fontSize="large" />}
        />
      ) : teams.length === 0 ? (
        <AppPageState
          description="Joined teams will appear here once you start collaborating with others."
          title="No teams yet"
          visual={<GroupRoundedIcon color="disabled" fontSize="large" />}
        />
      ) : (
        <AppStack spacing={0.5}>
          {teams.slice(0, 4).map((team, index) => (
            <AppBox
              component={RouterLink}
              key={team.id}
              sx={teamsCardItemLinkSx}
              to={AppRoutes.TEAM_DETAILS_PATH(team.id)}
            >
              <AppStack alignItems="center" direction="row" spacing={1.5}>
                <AppSurface sx={teamsCardAvatarSx(index)}>
                  {getInitials(team.name)}
                </AppSurface>
                <AppStack spacing={0.25}>
                  <AppTypography variant="subtitle2">{team.name}</AppTypography>
                  <AppTypography color="text.secondary" variant="caption">
                    Created {formatRelativeTime(team.createdAt)}
                  </AppTypography>
                </AppStack>
              </AppStack>
              <AppTypography color="primary.main" variant="caption">
                Open
              </AppTypography>
            </AppBox>
          ))}
        </AppStack>
      )}
    </AppSurface>
  </AppStack>
);
