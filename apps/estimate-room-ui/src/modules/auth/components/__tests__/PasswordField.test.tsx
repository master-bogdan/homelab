import userEvent from '@testing-library/user-event';

import { renderWithProviders, screen } from '@/test/test-utils';

import { PasswordField } from '../PasswordField';

describe('PasswordField', () => {
  it('toggles password visibility from the trailing action button', async () => {
    const user = userEvent.setup();

    renderWithProviders(
      <PasswordField
        label="Password"
        onChange={() => undefined}
        value="Password1!"
      />
    );

    expect(screen.getByLabelText('Password')).toHaveAttribute('type', 'password');

    await user.click(screen.getByRole('button', { name: 'Show password' }));

    expect(screen.getByLabelText('Password')).toHaveAttribute('type', 'text');

    await user.click(screen.getByRole('button', { name: 'Hide password' }));

    expect(screen.getByLabelText('Password')).toHaveAttribute('type', 'password');
  });
});
