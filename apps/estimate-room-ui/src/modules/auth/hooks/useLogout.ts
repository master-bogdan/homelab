import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { useAppDispatch } from '@/shared/store';
import { appRoutes } from '@/shared/constants/routes';

import { submitLogout } from '../store';

export const useLogout = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  const handleLogout = async () => {
    if (isLoggingOut) {
      return;
    }

    setIsLoggingOut(true);

    try {
      await dispatch(submitLogout()).unwrap();
    } finally {
      navigate(appRoutes.login, { replace: true });
      setIsLoggingOut(false);
    }
  };

  return {
    isLoggingOut,
    logout: handleLogout
  };
};
