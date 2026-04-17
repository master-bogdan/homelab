import { useEffect } from 'react';

import { useAppDispatch, useAppSelector } from '@/shared/store';
import type { AppDispatch } from '@/shared/store';
import { usePageTitle } from '@/shared/hooks';

import { fetchDashboardPage, selectDashboardPageState } from '../store';

let dashboardPageRequest: Promise<unknown> | null = null;

const getOrCreateDashboardPageRequest = (dispatch: AppDispatch) => {
  if (dashboardPageRequest) {
    return dashboardPageRequest;
  }

  const request = dispatch(fetchDashboardPage()).finally(() => {
    if (dashboardPageRequest === request) {
      dashboardPageRequest = null;
    }
  });

  dashboardPageRequest = request;

  return request;
};

export const useDashboardPage = () => {
  usePageTitle('Dashboard');
  const dispatch = useAppDispatch();
  const state = useAppSelector(selectDashboardPageState);

  useEffect(() => {
    getOrCreateDashboardPageRequest(dispatch);
  }, [dispatch]);

  return {
    ...state,
    retry: () => {
      dispatch(fetchDashboardPage());
    }
  };
};
