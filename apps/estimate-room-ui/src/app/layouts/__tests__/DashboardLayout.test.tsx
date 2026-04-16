import { Route, Routes, useLocation } from 'react-router-dom';
import userEvent from '@testing-library/user-event';

import { AuthStates } from '@/modules/auth';
import { AppRoutes } from '@/shared/constants/routes';
import { renderWithProviders, screen, waitFor } from '@/test/test-utils';

import { DashboardLayout } from '../DashboardLayout';

const LoginStateProbe = () => {
  const location = useLocation();
  const from = (
    location.state as { from?: { hash: string; pathname: string; search: string } } | null
  )?.from;

  return (
    <div>{from ? `${from.pathname}${from.search}${from.hash}` : 'missing redirect state'}</div>
  );
};

const createJsonResponse = (payload: unknown, status = 200) =>
  new Response(JSON.stringify(payload), {
    headers: {
      'content-type': 'application/json'
    },
    status
  });

describe('DashboardLayout', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(createJsonResponse([])));
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  it('renders route-aware metadata for protected pages', () => {
    renderWithProviders(
      <Routes>
        <Route element={<DashboardLayout />}>
          <Route element={<div>History page</div>} path={AppRoutes.HISTORY} />
        </Route>
      </Routes>,
      {
        preloadedState: {
          auth: {
            status: AuthStates.AUTHENTICATED,
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
          initialEntries: [AppRoutes.HISTORY]
        }
      }
    );

    expect(screen.getByRole('heading', { level: 6, name: 'History' })).toBeInTheDocument();
    expect(
      screen.getByText('Review completed sessions and archived room outcomes.')
    ).toBeInTheDocument();
    expect(screen.getByText('EstimateRoom Member')).toBeInTheDocument();
  });

  it('redirects unauthenticated users to login with the full requested URL in state', () => {
    renderWithProviders(
      <Routes>
        <Route element={<DashboardLayout />}>
          <Route element={<div>Private content</div>} path="rooms/:roomId" />
        </Route>
        <Route element={<LoginStateProbe />} path={AppRoutes.LOGIN} />
      </Routes>,
      {
        preloadedState: {
          auth: {
            status: AuthStates.UNAUTHENTICATED,
            user: null
          }
        },
        routerProps: {
          initialEntries: ['/rooms/room-123?tab=activity#voting']
        }
      }
    );

    expect(screen.getByText('/rooms/room-123?tab=activity#voting')).toBeInTheDocument();
  });

  it('renders the current user occupation when it exists', () => {
    renderWithProviders(
      <Routes>
        <Route element={<DashboardLayout />}>
          <Route element={<div>Dashboard page</div>} path={AppRoutes.DASHBOARD} />
        </Route>
      </Routes>,
      {
        preloadedState: {
          auth: {
            status: AuthStates.AUTHENTICATED,
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
          initialEntries: [AppRoutes.DASHBOARD]
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
          <Route element={<div>Dashboard page</div>} path={AppRoutes.DASHBOARD} />
        </Route>
      </Routes>,
      {
        preloadedState: {
          auth: {
            status: AuthStates.AUTHENTICATED,
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
          initialEntries: [AppRoutes.DASHBOARD]
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
          <Route element={<div>Dashboard page</div>} path={AppRoutes.DASHBOARD} />
        </Route>
      </Routes>,
      {
        preloadedState: {
          auth: {
            status: AuthStates.AUTHENTICATED,
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
          initialEntries: [AppRoutes.DASHBOARD]
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
