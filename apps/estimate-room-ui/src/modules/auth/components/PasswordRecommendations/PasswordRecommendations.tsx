import { AppBox, OverlineText } from '@/shared/components';

import { getPasswordRecommendations } from '../../utils';
import { PasswordRecommendationItem } from './PasswordRecommendationItem';
import {
  passwordRecommendationsGridSx,
  passwordRecommendationsRootSx,
  passwordRecommendationsTitleSx
} from './styles';

interface PasswordRecommendationsProps {
  readonly password: string;
}

export const PasswordRecommendations = ({
  password
}: PasswordRecommendationsProps) => {
  const recommendations = getPasswordRecommendations(password);

  return (
    <AppBox sx={passwordRecommendationsRootSx}>
      <OverlineText sx={passwordRecommendationsTitleSx}>Recommendations</OverlineText>
      <AppBox sx={passwordRecommendationsGridSx}>
        {recommendations.map((recommendation) => (
          <PasswordRecommendationItem
            key={recommendation.id}
            isMet={recommendation.isMet}
            label={recommendation.label}
          />
        ))}
      </AppBox>
    </AppBox>
  );
};
