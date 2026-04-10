import WorkspacePremiumRoundedIcon from '@mui/icons-material/WorkspacePremiumRounded';
import { LinearProgress, Paper, Stack, Typography } from '@mui/material';

import { AppButton, AppPageState } from '@/shared/ui';

import type { DashboardLedger } from '../../types';
import { getArchitectLevelLabel, getXpHint } from '../../utils';

import {
  ledgerCardBadgeSx,
  ledgerCardMetricsGridSx,
  ledgerCardProgressBarSx,
  ledgerCardRootSx
} from './styles';

export interface LedgerCardProps {
  readonly errorMessage: string | null;
  readonly ledger: DashboardLedger | null;
  readonly onRetry: () => void;
}

export const LedgerCard = ({ errorMessage, ledger, onRetry }: LedgerCardProps) => (
  <Stack spacing={1.5}>
    <Typography color="text.secondary" variant="overline">
      Architect Ledger
    </Typography>
    <Paper elevation={0} sx={ledgerCardRootSx(Boolean(ledger) && !errorMessage)}>
      {errorMessage || !ledger ? (
        <AppPageState
          action={
            <AppButton onClick={onRetry} variant="contained">
              Retry
            </AppButton>
          }
          description={errorMessage ?? 'No architect ledger data is available yet.'}
          title="Ledger unavailable"
          visual={<WorkspacePremiumRoundedIcon color="disabled" fontSize="large" />}
        />
      ) : (
        <Stack spacing={3}>
          <Stack alignItems="center" direction="row" spacing={2}>
            <Paper elevation={0} sx={ledgerCardBadgeSx}>
              <WorkspacePremiumRoundedIcon />
            </Paper>
            <Stack spacing={0.5}>
              <Typography variant="h4">{getArchitectLevelLabel(ledger.level)}</Typography>
              <Typography color="text.secondary" variant="caption">
                {ledger.currentLevelXp} / {ledger.nextLevelXp} XP
              </Typography>
            </Stack>
          </Stack>
          <Stack spacing={1}>
            <Stack alignItems="center" direction="row" justifyContent="space-between">
              <Typography color="text.secondary" variant="overline">
                Level Progress
              </Typography>
              <Typography color="primary.main" variant="caption">
                {ledger.xpProgressPercentage}%
              </Typography>
            </Stack>
            <LinearProgress sx={ledgerCardProgressBarSx} value={ledger.xpProgressPercentage} variant="determinate" />
            <Typography color="text.secondary" variant="caption">
              {getXpHint(ledger.currentLevelXp, ledger.nextLevelXp, ledger.level)}
            </Typography>
          </Stack>
          <Stack sx={ledgerCardMetricsGridSx}>
            <Stack spacing={0.5} textAlign="center">
              <Typography variant="h6">{ledger.tasksEstimated}</Typography>
              <Typography color="text.secondary" variant="caption">
                Tasks Estimated
              </Typography>
            </Stack>
            <Stack spacing={0.5} textAlign="center">
              <Typography variant="h6">{ledger.sessionsParticipated}</Typography>
              <Typography color="text.secondary" variant="caption">
                Sessions Joined
              </Typography>
            </Stack>
            <Stack spacing={0.5} textAlign="center">
              <Typography variant="h6">{ledger.sessionsAdmined}</Typography>
              <Typography color="text.secondary" variant="caption">
                Sessions Led
              </Typography>
            </Stack>
            <Stack spacing={0.5} textAlign="center">
              <Typography variant="h6">{ledger.achievements.length}</Typography>
              <Typography color="text.secondary" variant="caption">
                Achievements
              </Typography>
            </Stack>
          </Stack>
        </Stack>
      )}
    </Paper>
  </Stack>
);
