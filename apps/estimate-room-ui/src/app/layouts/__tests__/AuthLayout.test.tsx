import { Route, Routes } from 'react-router-dom';

import { AuthStates } from '@/modules/auth';
import { AppRoutes } from '@/shared/constants/routes';
import { renderWithProviders, screen } from '@/test/test-utils';

import { AuthLayout } from '../AuthLayout';

describe('AuthLayout', () => {
  it('redirects authenticated users away from login to the dashboard', () => {
    renderWithProviders(
      <Routes>
        <Route element={<AuthLayout />}>
          <Route element={<div>Login page</div>} path={AppRoutes.LOGIN} />
        </Route>
        <Route element={<div>Dashboard page</div>} path={AppRoutes.DASHBOARD} />
      </Routes>,
      {
        preloadedState: {
          auth: {
            oauthCallback: {
              errorMessage: null,
              redirectTo: null,
              requestKey: null,
              status: 'idle'
            },
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
          initialEntries: [AppRoutes.LOGIN]
        }
      }
    );

    expect(screen.getByText('Dashboard page')).toBeInTheDocument();
  });

  it('redirects authenticated users to the preserved protected route when present', () => {
    renderWithProviders(
      <Routes>
        <Route element={<AuthLayout />}>
          <Route element={<div>Login page</div>} path={AppRoutes.LOGIN} />
        </Route>
        <Route element={<div>Room page</div>} path={AppRoutes.ROOM_DETAILS} />
      </Routes>,
      {
        preloadedState: {
          auth: {
            oauthCallback: {
              errorMessage: null,
              redirectTo: null,
              requestKey: null,
              status: 'idle'
            },
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
          initialEntries: [
            {
              pathname: AppRoutes.LOGIN,
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
          <Route element={<div>Login page</div>} path={AppRoutes.LOGIN} />
        </Route>
      </Routes>,
      {
        preloadedState: {
          auth: {
            oauthCallback: {
              errorMessage: null,
              redirectTo: null,
              requestKey: null,
              status: 'idle'
            },
            status: AuthStates.UNAUTHENTICATED,
            user: null
          }
        },
        routerProps: {
          initialEntries: [AppRoutes.LOGIN]
        }
      }
    );

    expect(screen.getByText('Login page')).toBeInTheDocument();
  });
});
