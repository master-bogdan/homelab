import type { PropsWithChildren } from 'react';

import { AuthCard } from '../AuthCard';
import { authNarrowCardSx } from './styles';

export const AuthNarrowCard = ({ children }: PropsWithChildren) => (
  <AuthCard sx={authNarrowCardSx}>{children}</AuthCard>
);
