import { useEffect, useRef, useState } from 'react';
import { Link as RouterLink, useNavigate, useParams } from 'react-router-dom';

import { AppRoutes } from '@/app/router/routePaths';
import { useAppDispatch } from '@/shared/hooks';
import { submitJoinRoom } from '@/modules/dashboard/store';
import { AppButton, AppPageState } from '@/shared/components';

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

    const joinRoomRequest = dispatch(submitJoinRoom(token));

    joinRoomRequest.then((result: Awaited<typeof joinRoomRequest>) => {
      if (submitJoinRoom.fulfilled.match(result)) {
        navigate(AppRoutes.ROOM_DETAILS_PATH(result.payload.roomId), { replace: true });
        return;
      }

      if (submitJoinRoom.rejected.match(result)) {
        setError({
          message:
            typeof result.payload === 'string'
              ? result.payload
              : 'Invalid or expired room code. Please check and try again.',
          token
        });
      }
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
