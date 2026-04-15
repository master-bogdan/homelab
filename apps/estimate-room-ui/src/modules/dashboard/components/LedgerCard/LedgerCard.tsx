import WorkspacePremiumRoundedIcon from '@mui/icons-material/WorkspacePremiumRounded';

import {
  AppButton,
  AppPageState,
  AppProgress,
  AppStack,
  AppSurface,
  AppTypography
} from '@/shared/ui';

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
  <AppStack spacing={1.5}>
    <AppTypography color="text.secondary" variant="overline">
      Architect Ledger
    </AppTypography>
    <AppSurface sx={ledgerCardRootSx(Boolean(ledger) && !errorMessage)}>
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
        <AppStack spacing={3}>
          <AppStack alignItems="center" direction="row" spacing={2}>
            <AppSurface sx={ledgerCardBadgeSx}>
              <WorkspacePremiumRoundedIcon />
            </AppSurface>
            <AppStack spacing={0.5}>
              <AppTypography variant="h4">{getArchitectLevelLabel(ledger.level)}</AppTypography>
              <AppTypography color="text.secondary" variant="caption">
                {ledger.currentLevelXp} / {ledger.nextLevelXp} XP
              </AppTypography>
            </AppStack>
          </AppStack>
          <AppStack spacing={1}>
            <AppStack alignItems="center" direction="row" justifyContent="space-between">
              <AppTypography color="text.secondary" variant="overline">
                Level Progress
              </AppTypography>
              <AppTypography color="primary.main" variant="caption">
                {ledger.xpProgressPercentage}%
              </AppTypography>
            </AppStack>
            <AppProgress kind="linear" sx={ledgerCardProgressBarSx} value={ledger.xpProgressPercentage} variant="determinate" />
            <AppTypography color="text.secondary" variant="caption">
              {getXpHint(ledger.currentLevelXp, ledger.nextLevelXp, ledger.level)}
            </AppTypography>
          </AppStack>
          <AppStack sx={ledgerCardMetricsGridSx}>
            <AppStack spacing={0.5} textAlign="center">
              <AppTypography variant="h6">{ledger.tasksEstimated}</AppTypography>
              <AppTypography color="text.secondary" variant="caption">
                Tasks Estimated
              </AppTypography>
            </AppStack>
            <AppStack spacing={0.5} textAlign="center">
              <AppTypography variant="h6">{ledger.sessionsParticipated}</AppTypography>
              <AppTypography color="text.secondary" variant="caption">
                Sessions Joined
              </AppTypography>
            </AppStack>
            <AppStack spacing={0.5} textAlign="center">
              <AppTypography variant="h6">{ledger.sessionsAdmined}</AppTypography>
              <AppTypography color="text.secondary" variant="caption">
                Sessions Led
              </AppTypography>
            </AppStack>
            <AppStack spacing={0.5} textAlign="center">
              <AppTypography variant="h6">{ledger.achievements.length}</AppTypography>
              <AppTypography color="text.secondary" variant="caption">
                Achievements
              </AppTypography>
            </AppStack>
          </AppStack>
        </AppStack>
      )}
    </AppSurface>
  </AppStack>
);
