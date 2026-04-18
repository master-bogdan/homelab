import { renderWithProviders, screen } from '@/test/test-utils';

import { AppDialog } from '../AppDialog';

describe('AppDialog', () => {
  it('renders a titled dialog with shared cancel and confirm actions', () => {
    renderWithProviders(
      <AppDialog
        confirmLabel="Confirm"
        onClose={() => undefined}
        open
        title="Shared Dialog"
      >
        Dialog content
      </AppDialog>
    );

    expect(screen.getByRole('dialog', { name: 'Shared Dialog' })).toBeInTheDocument();
    expect(screen.getByText('Dialog content')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Confirm' })).toBeInTheDocument();
  });

  it('renders content without a title and supports hiding footer buttons', () => {
    renderWithProviders(
      <AppDialog
        hideCancelButton
        hideConfirmButton
        onClose={() => undefined}
        open
      >
        Untitled dialog
      </AppDialog>
    );

    expect(screen.getByRole('dialog')).toBeInTheDocument();
    expect(screen.getByText('Untitled dialog')).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Cancel' })).not.toBeInTheDocument();
  });
});
