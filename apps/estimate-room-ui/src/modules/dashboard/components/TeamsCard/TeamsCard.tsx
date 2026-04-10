import GroupRoundedIcon from '@mui/icons-material/GroupRounded';
import { Box, Paper, Stack, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { AppButton, AppPageState } from '@/shared/ui';

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
  <Stack spacing={1.5}>
    <Typography color="text.secondary" variant="overline">
      Your Teams
    </Typography>
    <Paper elevation={0} sx={teamsCardRootSx(teams.length === 0 || Boolean(errorMessage))}>
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
        <Stack spacing={0.5}>
          {teams.slice(0, 4).map((team, index) => (
            <Box
              component={RouterLink}
              key={team.id}
              sx={teamsCardItemLinkSx}
              to={appRoutes.teamDetailsPath(team.id)}
            >
              <Stack alignItems="center" direction="row" spacing={1.5}>
                <Paper elevation={0} sx={teamsCardAvatarSx(index)}>
                  {getInitials(team.name)}
                </Paper>
                <Stack spacing={0.25}>
                  <Typography variant="subtitle2">{team.name}</Typography>
                  <Typography color="text.secondary" variant="caption">
                    Created {formatRelativeTime(team.createdAt)}
                  </Typography>
                </Stack>
              </Stack>
              <Typography color="primary.main" variant="caption">
                Open
              </Typography>
            </Box>
          ))}
        </Stack>
      )}
    </Paper>
  </Stack>
);
