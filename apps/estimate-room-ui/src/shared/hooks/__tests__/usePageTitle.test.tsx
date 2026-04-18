import { renderHook } from '@testing-library/react';

import { AppConfig } from '@/config';

import { usePageTitle } from '../usePageTitle';

describe('usePageTitle', () => {
  it('sets and restores the document title', () => {
    const originalTitle = document.title;
    const { unmount } = renderHook(() => usePageTitle('Dashboard'));

    expect(document.title).toBe(`Dashboard | ${AppConfig.APP_NAME}`);

    unmount();

    expect(document.title).toBe(originalTitle);
  });
});
