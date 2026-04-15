import { useEffect } from 'react';

import { appConfig } from '@/config';

export const usePageTitle = (pageTitle: string) => {
  useEffect(() => {
    const previousTitle = document.title;
    document.title = `${pageTitle} | ${appConfig.appName}`;

    return () => {
      document.title = previousTitle;
    };
  }, [pageTitle]);
};
