import { Route, Routes } from 'react-router-dom';

import { AUTH_STATUSES } from '@/modules/auth/types';
import { appRoutes } from '@/shared/constants/routes';
import { renderWithProviders, screen } from '@/test/test-utils';

import { AuthLayout } from '../AuthLayout';

describe('AuthLayout', () => {
  it('redirects authenticated users away from login to the dashboard', () => {
    renderWithProviders(
      <Routes>
        <Route element={<AuthLayout />}>
          <Route element={<div>Login page</div>} path={appRoutes.login} />
        </Route>
        <Route element={<div>Dashboard page</div>} path={appRoutes.dashboard} />
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
          initialEntries: [appRoutes.login]
        }
      }
    );

    expect(screen.getByText('Dashboard page')).toBeInTheDocument();
  });

  it('redirects authenticated users to the preserved protected route when present', () => {
    renderWithProviders(
      <Routes>
        <Route element={<AuthLayout />}>
          <Route element={<div>Login page</div>} path={appRoutes.login} />
        </Route>
        <Route element={<div>Room page</div>} path={appRoutes.roomDetails} />
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
          initialEntries: [
            {
              pathname: appRoutes.login,
              state: {
                from: {
                  hash: '#voting',
                  pathname: '/rooms/room-123',
                  search: '?tab=activity'
                }
              }
            }
          ]
        }
      }
    );

    expect(screen.getByText('Room page')).toBeInTheDocument();
  });

  it('keeps auth pages visible for unauthenticated users', () => {
    renderWithProviders(
      <Routes>
        <Route element={<AuthLayout />}>
          <Route element={<div>Login page</div>} path={appRoutes.login} />
        </Route>
      </Routes>,
      {
        preloadedState: {
          auth: {
            status: AUTH_STATUSES.UNAUTHENTICATED,
            user: null
          }
        },
        routerProps: {
          initialEntries: [appRoutes.login]
        }
      }
    );

    expect(screen.getByText('Login page')).toBeInTheDocument();
  });
});
