import { useEffect } from 'react';

import { AppConfig } from '@/config';

export const usePageTitle = (pageTitle: string) => {
  useEffect(() => {
    const previousTitle = document.title;
    document.title = `${pageTitle} | ${AppConfig.APP_NAME}`;

    return () => {
      document.title = previousTitle;
    };
  }, [pageTitle]);
};
