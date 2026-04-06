import type { AuthUser } from '@/shared/types';

import type { ProfileSummary } from '../types';

export const profileService = {
  getProfileSummary: (user: AuthUser | null): ProfileSummary => ({
    displayName: user?.displayName ?? 'Awaiting backend session',
    email: user?.email ?? 'No authenticated user yet',
    roleLabel: user ? 'session-authenticated' : 'pending',
    teamCount: 0
  })
};
