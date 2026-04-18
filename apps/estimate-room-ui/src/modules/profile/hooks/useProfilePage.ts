import { useAppSelector } from '@/shared/hooks';
import { selectAuthUser } from '@/modules/auth';
import { usePageTitle } from '@/shared/hooks';

import { profileService } from '../api/profileApi';

export const useProfilePage = () => {
  const user = useAppSelector(selectAuthUser);

  usePageTitle('Profile');

  return {
    profile: profileService.getProfileSummary(user)
  };
};
