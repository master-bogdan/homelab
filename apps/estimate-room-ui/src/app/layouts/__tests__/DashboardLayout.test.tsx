import { Route, Routes } from 'react-router-dom';
import userEvent from '@testing-library/user-event';

import { AUTH_STATUSES } from '@/modules/auth/types';
import { dashboardService } from '@/modules/dashboard/services/dashboardService';
import { appRoutes } from '@/shared/constants/routes';
import { renderWithProviders, screen, waitFor } from '@/test/test-utils';

import { DashboardLayout } from '../DashboardLayout';

describe('DashboardLayout', () => {
  beforeEach(() => {
    vi.spyOn(dashboardService, 'fetchTeams').mockResolvedValue([]);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('renders route-aware metadata for protected pages', () => {
    renderWithProviders(
      <Routes>
        <Route element={<DashboardLayout />}>
          <Route element={<div>History page</div>} path={appRoutes.history} />
        </Route>
      </Routes>,
      {
        preloadedState: {
          auth: {
            status: AUTH_STATUSES.AUTHENTICATED,
            user: {
              avatarUrl: null,
              displayName: 'Alex Architect',
              email: 'alex@example.com',
              id: 'user-1',
              occupation: null,
              organization: null
            }
          }
        },
        routerProps: {
          initialEntries: [appRoutes.history]
        }
      }
    );

    expect(screen.getByRole('heading', { level: 6, name: 'History' })).toBeInTheDocument();
    expect(
      screen.getByText('Review completed sessions and archived room outcomes.')
    ).toBeInTheDocument();
    expect(screen.getByText('EstimateRoom Member')).toBeInTheDocument();
  });

  it('renders the current user occupation when it exists', () => {
    renderWithProviders(
      <Routes>
        <Route element={<DashboardLayout />}>
          <Route element={<div>Dashboard page</div>} path={appRoutes.dashboard} />
        </Route>
      </Routes>,
      {
        preloadedState: {
          auth: {
            status: AUTH_STATUSES.AUTHENTICATED,
            user: {
              avatarUrl: null,
              displayName: 'Alex Architect',
              email: 'alex@example.com',
              id: 'user-1',
              occupation: 'Technical Architect',
              organization: null
            }
          }
        },
        routerProps: {
          initialEntries: [appRoutes.dashboard]
        }
      }
    );

    expect(screen.getByText('Technical Architect')).toBeInTheDocument();
  });

  it('shows the user menu summary and actions', async () => {
    const user = userEvent.setup();

    renderWithProviders(
      <Routes>
        <Route element={<DashboardLayout />}>
          <Route element={<div>Dashboard page</div>} path={appRoutes.dashboard} />
        </Route>
      </Routes>,
      {
        preloadedState: {
          auth: {
            status: AUTH_STATUSES.AUTHENTICATED,
            user: {
              avatarUrl: null,
              displayName: 'Alex Architect',
              email: 'alex@example.com',
              id: 'user-1',
              occupation: null,
              organization: null
            }
          }
        },
        routerProps: {
          initialEntries: [appRoutes.dashboard]
        }
      }
    );

    await user.click(screen.getByLabelText('Open user menu'));

    expect(screen.getByText('alex@example.com')).toBeInTheDocument();
    expect(screen.getByRole('menuitem', { name: 'Profile' })).toBeInTheDocument();
    expect(screen.getByRole('menuitem', { name: 'Settings' })).toBeInTheDocument();
    expect(screen.getByRole('menuitem', { name: 'Log out' })).toBeInTheDocument();
  });

  it('opens the shared create-room and join-room dialogs from the header', async () => {
    const user = userEvent.setup();

    renderWithProviders(
      <Routes>
        <Route element={<DashboardLayout />}>
          <Route element={<div>Dashboard page</div>} path={appRoutes.dashboard} />
        </Route>
      </Routes>,
      {
        preloadedState: {
          auth: {
            status: AUTH_STATUSES.AUTHENTICATED,
            user: {
              avatarUrl: null,
              displayName: 'Alex Architect',
              email: 'alex@example.com',
              id: 'user-1',
              occupation: null,
              organization: null
            }
          }
        },
        routerProps: {
          initialEntries: [appRoutes.dashboard]
        }
      }
    );

    await user.click(screen.getByRole('button', { name: 'Create room' }));
    expect(await screen.findByText('Create New Room')).toBeInTheDocument();

    await user.click(screen.getByRole('button', { name: 'Cancel' }));
    await waitFor(() => {
      expect(screen.queryByText('Create New Room')).not.toBeInTheDocument();
    });
    await user.click(screen.getByRole('button', { name: 'Join room' }));

    expect(await screen.findByRole('dialog', { name: 'Join Room' })).toBeInTheDocument();
  });
});
