import { expect, test } from '@playwright/test';

import { mockAuthenticatedSession, mockDashboardApi } from './fixtures/api';

test('renders the authenticated dashboard with mocked backend data', async ({ page }) => {
  await mockAuthenticatedSession(page);
  await mockDashboardApi(page);

  await page.goto('/dashboard');

  await expect(page.getByText('No active rooms yet')).toBeVisible();
  await expect(page.getByText('No teams yet')).toBeVisible();
  await expect(page.getByText('Architect Ledger')).toBeVisible();
});
