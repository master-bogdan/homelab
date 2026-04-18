import { useParams } from 'react-router-dom';

import { usePageTitle } from '@/shared/hooks';

import { roomsService } from '../api/roomsApi';

export const useRoomDetailsPage = () => {
  const { id = '' } = useParams();
  const room = roomsService.getRoomPreview(id);

  usePageTitle(room?.name ?? 'Room Details');

  return {
    room,
    roomId: id
  };
};
