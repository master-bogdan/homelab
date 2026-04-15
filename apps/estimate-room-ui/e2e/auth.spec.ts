import { expect, test } from '@playwright/test';

import { mockAuthMutations, mockGuestSession } from './fixtures/api';

test.beforeEach(async ({ page }) => {
  await mockGuestSession(page);
  await mockAuthMutations(page);
});

test('shows validation feedback on the login page', async ({ page }) => {
  await page.goto('/login');

  await page.getByPlaceholder('name@company.com').fill('not-an-email');
  await page.getByPlaceholder('••••••••').fill('short');
  await page.getByRole('button', { name: 'Continue with GitHub' }).focus();

  await expect(page.getByText('Enter a valid email address.')).toBeVisible();
  await expect(page.getByText('Password must be at least 8 characters.')).toBeVisible();
});

test('navigates between auth pages', async ({ page }) => {
  await page.goto('/login');

  await page.getByRole('link', { name: 'Register now' }).click();
  await expect(page.getByRole('heading', { name: 'Create Workspace' })).toBeVisible();

  await page.getByRole('link', { name: 'Sign In' }).click();
  await expect(page.getByRole('heading', { name: 'Welcome Back' })).toBeVisible();

  await page.getByRole('link', { name: 'Forgot?' }).click();
  await expect(page).toHaveURL(/\/forgot-password$/);
});

test('submits the forgot password form', async ({ page }) => {
  await page.goto('/forgot-password');

  await page.getByPlaceholder('name@company.com').fill('ada@example.com');
  await page.getByRole('button', { name: /send reset link/i }).click();

  await expect(page.getByText(/check your inbox/i)).toBeVisible();
});
