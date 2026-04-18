import type {
  PasswordRecommendationRuleIds,
  PasswordRecommendationRules
} from '../constants/passwordRules';

export type PasswordRecommendationRuleId =
  (typeof PasswordRecommendationRuleIds)[keyof typeof PasswordRecommendationRuleIds];

export type PasswordRecommendationRule = (typeof PasswordRecommendationRules)[number];

export interface PasswordRecommendation {
  readonly id: PasswordRecommendationRuleId;
  readonly isMet: boolean;
  readonly label: string;
}
