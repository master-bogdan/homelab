import { renderWithProviders, screen } from '@/test/test-utils';

import { DashboardPage } from '../DashboardPage';

describe('DashboardPage', () => {
  it('renders the dashboard heading and primary action', () => {
    renderWithProviders(<DashboardPage />);

    expect(screen.getByRole('heading', { name: 'Dashboard' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Create room' })).toBeInTheDocument();
  });
});
