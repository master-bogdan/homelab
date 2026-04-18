import {
  isEmailAlreadyInUseError,
  isInvalidCredentialsError,
  resolveApiErrorMessage
} from '../errorMessages';

describe('auth error message utilities', () => {
  it('resolves API payload messages by priority', () => {
    expect(
      resolveApiErrorMessage(
        {
          detail: 'Detail message',
          message: 'Message text',
          status: 400,
          title: 'Title text'
        },
        'Fallback'
      )
    ).toBe('Detail message');
  });

  it('resolves RTK Query payload and transport messages', () => {
    expect(
      resolveApiErrorMessage(
        {
          data: {
            message: 'RTK message'
          },
          status: 400
        },
        'Fallback'
      )
    ).toBe('RTK message');

    expect(
      resolveApiErrorMessage(
        {
          error: 'Network Error',
          status: 'FETCH_ERROR'
        },
        'Fallback'
      )
    ).toBe('Network Error');
  });

  it('falls back when no readable message exists', () => {
    expect(resolveApiErrorMessage({ status: 500 }, 'Fallback')).toBe('Fallback');
  });

  it('detects known auth error messages across error shapes', () => {
    expect(isInvalidCredentialsError(new Error('Invalid credentials'))).toBe(true);
    expect(isEmailAlreadyInUseError({
      data: {
        detail: 'Email already in use'
      },
      status: 409
    })).toBe(true);
  });
});
