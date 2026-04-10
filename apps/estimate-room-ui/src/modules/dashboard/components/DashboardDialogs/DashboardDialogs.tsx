import { useCreateRoomDialog } from '../../hooks/useCreateRoomDialog';
import { useJoinRoomDialog } from '../../hooks/useJoinRoomDialog';
import { CreateRoomDialog } from '../CreateRoomDialog';
import { CreateRoomSuccessDialog } from '../CreateRoomSuccessDialog';
import { JoinRoomDialog } from '../JoinRoomDialog';

export const DashboardDialogs = () => {
  const createRoomDialog = useCreateRoomDialog();
  const joinRoomDialog = useJoinRoomDialog();

  return (
    <>
      <CreateRoomDialog
        form={createRoomDialog.form}
        isLoadingTeams={createRoomDialog.isLoadingTeams}
        onClose={createRoomDialog.close}
        onSubmit={createRoomDialog.onSubmit}
        open={createRoomDialog.isOpen}
        submitErrorMessage={createRoomDialog.submitErrorMessage}
        teamErrorMessage={createRoomDialog.teamErrorMessage}
        teamOptions={createRoomDialog.teamOptions}
      />
      <CreateRoomSuccessDialog
        onClose={createRoomDialog.closeResult}
        onOpenRoom={createRoomDialog.openCreatedRoom}
        result={createRoomDialog.result}
      />
      <JoinRoomDialog
        errorMessage={joinRoomDialog.errorMessage}
        form={joinRoomDialog.form}
        onClose={joinRoomDialog.close}
        onSubmit={joinRoomDialog.onSubmit}
        open={joinRoomDialog.isOpen}
      />
    </>
  );
};
