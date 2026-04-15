import { renderWithProviders, screen } from '@/test/test-utils';

import { AppProgress } from '../AppProgress';

describe('AppProgress', () => {
  it('renders circular progress by default', () => {
    renderWithProviders(<AppProgress aria-label="Loading" />);

    expect(screen.getByRole('progressbar', { name: 'Loading' })).toBeInTheDocument();
  });

  it('renders linear progress when requested', () => {
    renderWithProviders(
      <AppProgress
        aria-label="Completion"
        kind="linear"
        value={40}
        variant="determinate"
      />
    );

    expect(screen.getByRole('progressbar', { name: 'Completion' })).toHaveAttribute(
      'aria-valuenow',
      '40'
    );
  });
});
