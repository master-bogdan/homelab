import { renderHook } from '@testing-library/react';

import { appConfig } from '@/config';

import { usePageTitle } from '../usePageTitle';

describe('usePageTitle', () => {
  it('sets and restores the document title', () => {
    const originalTitle = document.title;
    const { unmount } = renderHook(() => usePageTitle('Dashboard'));

    expect(document.title).toBe(`Dashboard | ${appConfig.appName}`);

    unmount();

    expect(document.title).toBe(originalTitle);
  });
});
