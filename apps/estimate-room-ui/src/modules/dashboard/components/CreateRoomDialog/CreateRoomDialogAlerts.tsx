import { AppAlert } from '@/shared/components';

interface CreateRoomDialogAlertsProps {
  readonly submitErrorMessage: string | null;
  readonly teamErrorMessage: string | null;
}

export const CreateRoomDialogAlerts = ({
  submitErrorMessage,
  teamErrorMessage
}: CreateRoomDialogAlertsProps) => (
  <>
    {submitErrorMessage ? <AppAlert severity="error">{submitErrorMessage}</AppAlert> : null}
    {teamErrorMessage ? (
      <AppAlert severity="warning">
        {teamErrorMessage} You can still create a room without linking a team.
      </AppAlert>
    ) : null}
  </>
);
