import { ForgotPasswordPage, LoginPage, RegisterPage } from '@/app/pages';
import { renderWithProviders, screen } from '@/test/test-utils';

describe('auth form actions', () => {
  it('disables register submit while the form is invalid', () => {
    renderWithProviders(<RegisterPage />);

    expect(screen.getByRole('button', { name: 'Initialize Account' })).toBeDisabled();
    expect(screen.getByRole('button', { name: 'Sign up with GitHub' })).toBeEnabled();
  });

  it('disables login submit while the form is invalid', () => {
    renderWithProviders(<LoginPage />);

    expect(screen.getByRole('button', { name: 'Sign In' })).toBeDisabled();
    expect(screen.getByRole('button', { name: 'Continue with GitHub' })).toBeEnabled();
  });

  it('disables forgot-password submit while the form is invalid', () => {
    renderWithProviders(<ForgotPasswordPage />);

    expect(screen.getByRole('button', { name: 'Send Reset Link' })).toBeDisabled();
  });
});
