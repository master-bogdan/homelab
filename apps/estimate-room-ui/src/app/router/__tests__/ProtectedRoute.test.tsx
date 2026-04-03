import { Route, Routes, useLocation } from 'react-router-dom';

import { appRoutes } from '@/shared/constants/routes';
import { renderWithProviders, screen } from '@/test/test-utils';

import { ProtectedRoute } from '../ProtectedRoute';

const LoginStateProbe = () => {
  const location = useLocation();
  const from = (location.state as { from?: { hash: string; pathname: string; search: string } } | null)
    ?.from;

  return <div>{from ? `${from.pathname}${from.search}${from.hash}` : 'missing redirect state'}</div>;
};

describe('ProtectedRoute', () => {
  it('preserves the full requested URL in redirect state', () => {
    renderWithProviders(
      <Routes>
        <Route element={<ProtectedRoute />}>
          <Route element={<div>Private content</div>} path="rooms/:roomId" />
        </Route>
        <Route element={<LoginStateProbe />} path={appRoutes.login} />
      </Routes>,
      {
        preloadedState: {
          auth: {
            status: 'unauthenticated',
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
});
