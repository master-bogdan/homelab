import { PasswordRecommendationRules } from '../constants';
import type { PasswordRecommendation } from '../types';

export const getPasswordRecommendations = (password: string): PasswordRecommendation[] =>
  PasswordRecommendationRules.map((rule) => ({
    id: rule.id,
    isMet: rule.test(password),
    label: rule.label
  }));
