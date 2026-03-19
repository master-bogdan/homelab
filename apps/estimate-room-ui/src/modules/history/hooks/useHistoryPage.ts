import { usePageTitle } from '@/shared/hooks';

import { historyService } from '../services/historyService';

export const useHistoryPage = () => {
  usePageTitle('History');

  return {
    entries: historyService.getHistoryEntries()
  };
};
