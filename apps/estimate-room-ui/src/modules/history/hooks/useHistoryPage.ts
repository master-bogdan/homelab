import { usePageTitle } from '@/shared/hooks';

import { historyService } from '../api/historyApi';

export const useHistoryPage = () => {
  usePageTitle('History');

  return {
    entries: historyService.getHistoryEntries()
  };
};
