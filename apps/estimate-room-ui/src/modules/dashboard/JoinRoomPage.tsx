import { useEffect, useRef, useState } from 'react';
import { Link as RouterLink, useNavigate, useParams } from 'react-router-dom';

import { AppRoutes } from '@/shared/constants/routes';
import { useAppDispatch } from '@/shared/store';
import { AppButton, AppPageState } from '@/shared/ui';

import { submitJoinRoom } from './store';

interface JoinRoomErrorState {
  readonly message: string;
  readonly token: string;
}

export const JoinRoomPage = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { token } = useParams<{ token: string }>();
  const [error, setError] = useState<JoinRoomErrorState | null>(null);
  const submittedTokenRef = useRef<string | null>(null);

  useEffect(() => {
    if (!token || submittedTokenRef.current === token) {
      return;
    }

    submittedTokenRef.current = token;

    void dispatch(submitJoinRoom(token))
      .unwrap()
      .then((result) => {
        navigate(AppRoutes.ROOM_DETAILS_PATH(result.roomId), { replace: true });
      })
      .catch((error: unknown) => {
        setError({
          message:
            typeof error === 'string'
              ? error
              : 'Invalid or expired room code. Please check and try again.',
          token
        });
      });
  }, [dispatch, navigate, token]);

  if (!token) {
    return (
      <AppPageState
        action={
          <AppButton component={RouterLink} to={AppRoutes.DASHBOARD} variant="contained">
            Back to Dashboard
          </AppButton>
        }
        description="Ask the room owner for a fresh invitation link."
        title="Invite link is missing"
      />
    );
  }

  if (error?.token === token) {
    return (
      <AppPageState
        action={
          <AppButton component={RouterLink} to={AppRoutes.DASHBOARD} variant="contained">
            Back to Dashboard
          </AppButton>
        }
        description={error.message}
        title="Unable to join room"
      />
    );
  }

  return (
    <AppPageState
      description="Accepting your room invitation."
      isLoading
      title="Joining room"
    />
  );
};
