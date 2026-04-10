import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { useAppDispatch } from '@/app/store/hooks';
import { appRoutes } from '@/shared/constants/routes';
import { api } from '@/shared/api';

import { authService } from '../services';
import { clearSession } from '../store';

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
      await authService.logout(dispatch);
    } finally {
      dispatch(api.util.resetApiState());
      dispatch(clearSession());
      navigate(appRoutes.login, { replace: true });
      setIsLoggingOut(false);
    }
  };

  return {
    isLoggingOut,
    logout: handleLogout
  };
};
