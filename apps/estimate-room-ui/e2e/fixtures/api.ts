import type { Page, Route } from '@playwright/test';

const json = async (route: Route, body: unknown, status = 200) => {
  await route.fulfill({
    body: JSON.stringify(body),
    contentType: 'application/json',
    status
  });
};

export const testUser = {
  avatarUrl: null,
  displayName: 'Ada Lovelace',
  email: 'ada@example.com',
  id: 'user-1',
  occupation: 'Engineer',
  organization: 'Analytical Engines'
};

export const mockGuestSession = async (page: Page) => {
  await page.route('**/api/v1/oauth2/token', (route) =>
    json(route, { message: 'No refresh token.' }, 401)
  );
  await page.route('**/api/v1/auth/logout', (route) =>
    json(route, { loggedOut: true })
  );
  await page.route('**/api/v1/auth/session', (route) =>
    json(route, { authenticated: false, user: null })
  );
};

export const mockAuthenticatedSession = async (page: Page) => {
  await page.addInitScript(() => {
    window.localStorage.setItem('estimate-room.auth.access-token', 'e2e-access-token');
  });
  await page.route('**/api/v1/auth/session', (route) =>
    json(route, { authenticated: true, user: testUser })
  );
};

export const mockAuthMutations = async (page: Page) => {
  await page.route('**/api/v1/auth/login', (route) =>
    json(route, { authenticated: true, user: testUser })
  );
  await page.route('**/api/v1/auth/register', (route) =>
    json(route, { authenticated: true, user: testUser })
  );
  await page.route('**/api/v1/auth/forgot-password', (route) =>
    json(route, { submitted: true })
  );
};

export const mockDashboardApi = async (page: Page) => {
  await page.route('**/api/v1/history/me/sessions**', (route) =>
    json(route, {
      items: [],
      page: 1,
      pageSize: 20,
      total: 0
    })
  );
  await page.route('**/api/v1/teams', (route) => json(route, []));
  await page.route('**/api/v1/gamification/me', (route) =>
    json(route, {
      achievements: [],
      stats: {
        level: 1,
        nextLevelXp: 100,
        sessionsAdmined: 0,
        sessionsParticipated: 0,
        tasksEstimated: 0,
        xp: 0
      }
    })
  );
};
