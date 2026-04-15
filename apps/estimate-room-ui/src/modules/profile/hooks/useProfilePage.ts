import { useAppSelector } from '@/shared/store';
import { selectAuthUser } from '@/modules/auth';
import { usePageTitle } from '@/shared/hooks';

import { profileService } from '../services/profileService';

export const useProfilePage = () => {
  const user = useAppSelector(selectAuthUser);

  usePageTitle('Profile');

  return {
    profile: profileService.getProfileSummary(user)
  };
};
