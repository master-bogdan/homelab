import { renderWithProviders, screen } from '@/test/test-utils';

import { AppButton } from '../../AppButton';
import { AppPageState } from '../AppPageState';

describe('AppPageState', () => {
  it('renders a loading state with copy and fallback spinner', () => {
    renderWithProviders(
      <AppPageState
        description="Checking your current session."
        isLoading
        title="Loading workspace"
      />
    );

    expect(screen.getByRole('progressbar')).toBeInTheDocument();
    expect(screen.getByText('Loading workspace')).toBeInTheDocument();
    expect(screen.getByText('Checking your current session.')).toBeInTheDocument();
  });

  it('renders an action when provided', () => {
    renderWithProviders(
      <AppPageState
        action={<AppButton variant="contained">Retry</AppButton>}
        description="Try the last action again."
        title="Something went wrong"
      />
    );

    expect(screen.getByRole('button', { name: 'Retry' })).toBeInTheDocument();
  });
});
