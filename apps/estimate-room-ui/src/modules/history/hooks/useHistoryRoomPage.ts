import { useParams } from 'react-router-dom';

import { usePageTitle } from '@/shared/hooks';

import { historyService } from '../api/historyApi';

export const useHistoryRoomPage = () => {
  const { id = '' } = useParams();

  usePageTitle(`Room History ${id}`);

  return {
    entries: historyService.getHistoryForRoom(id),
    roomId: id
  };
};
