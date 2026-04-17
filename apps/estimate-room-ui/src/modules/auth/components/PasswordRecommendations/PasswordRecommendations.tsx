import CheckCircleRoundedIcon from '@mui/icons-material/CheckCircleRounded';
import RadioButtonUncheckedRoundedIcon from '@mui/icons-material/RadioButtonUncheckedRounded';

import { AppBox, AppStack, AppTypography, OverlineText } from '@/shared/ui';

import {
  passwordRecommendationsGridSx,
  passwordRecommendationsRootSx,
  passwordRecommendationsTitleSx
} from './styles';

const getRecommendations = (password: string) => [
  {
    isMet: password.length >= 8,
    label: '8+ characters'
  },
  {
    isMet: /[0-9]/u.test(password),
    label: 'One number'
  },
  {
    isMet: /[A-Z]/u.test(password),
    label: 'Uppercase letter'
  },
  {
    isMet: /[^A-Za-z0-9]/u.test(password),
    label: 'Special symbol'
  }
];

interface PasswordRecommendationItemProps {
  readonly isMet: boolean;
  readonly label: string;
}

const PasswordRecommendationItem = ({
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

interface PasswordRecommendationsProps {
  readonly password: string;
}

export const PasswordRecommendations = ({
  password
}: PasswordRecommendationsProps) => {
  const recommendations = getRecommendations(password);

  return (
    <AppBox sx={passwordRecommendationsRootSx}>
      <OverlineText sx={passwordRecommendationsTitleSx}>Recommendations</OverlineText>
      <AppBox sx={passwordRecommendationsGridSx}>
        {recommendations.map((recommendation) => (
          <PasswordRecommendationItem
            key={recommendation.label}
            isMet={recommendation.isMet}
            label={recommendation.label}
          />
        ))}
      </AppBox>
    </AppBox>
  );
};
