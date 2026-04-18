import { renderWithProviders, screen } from '@/test/test-utils';

import { AppButton } from '../AppButton';

describe('AppButton', () => {
  it('renders the loading state with disabled interaction', () => {
    renderWithProviders(
      <AppButton loading loadingText="Saving changes" variant="contained">
        Save
      </AppButton>
    );

    expect(
      screen.getByRole('button', { name: 'Saving changes' })
    ).toBeDisabled();
    expect(screen.getByRole('button', { name: 'Saving changes' })).toHaveAttribute(
      'aria-busy',
      'true'
    );
    expect(screen.getByRole('progressbar')).toBeInTheDocument();
  });
});
