import CheckCircleRoundedIcon from '@mui/icons-material/CheckCircleRounded';
import RadioButtonUncheckedRoundedIcon from '@mui/icons-material/RadioButtonUncheckedRounded';

import { AppStack, AppTypography } from '@/shared/components';

interface PasswordRecommendationItemProps {
  readonly isMet: boolean;
  readonly label: string;
}

export const PasswordRecommendationItem = ({
  isMet,
  label
}: PasswordRecommendationItemProps) => {
  const Icon = isMet ? CheckCircleRoundedIcon : RadioButtonUncheckedRoundedIcon;

  return (
    <AppStack alignItems="center" direction="row" spacing={1}>
      <Icon color={isMet ? 'primary' : 'disabled'} fontSize="small" />
      <AppTypography color="text.secondary" variant="caption">
        {label}
      </AppTypography>
    </AppStack>
  );
};
