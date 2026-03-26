import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { usePageTitle } from '@/shared/hooks';

import { roomsService } from '../services/roomsService';
import type { NewRoomFormValues } from '../types';

export const newRoomDefaultValues: NewRoomFormValues = {
  height: 2.7,
  length: 6,
  name: '',
  teamId: '',
  width: 4
};

export const useNewRoomForm = () => {
  usePageTitle('New Room');

  const form = useForm<NewRoomFormValues>({
    defaultValues: newRoomDefaultValues,
    mode: 'onChange',
    reValidateMode: 'onChange'
  });
  const [submitMessage, setSubmitMessage] = useState<string | null>(null);

  const onSubmit = form.handleSubmit(async (values) => {
    const room = await roomsService.createRoom(values);

    setSubmitMessage(`Room "${room.name}" is ready for backend API integration.`);
    form.reset(newRoomDefaultValues);
  });

  return {
    form,
    onSubmit,
    submitMessage
  };
};
