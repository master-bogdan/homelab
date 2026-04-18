import CheckCircleRoundedIcon from '@mui/icons-material/CheckCircleRounded';
import ContentCopyRoundedIcon from '@mui/icons-material/ContentCopyRounded';
import WarningAmberRoundedIcon from '@mui/icons-material/WarningAmberRounded';
import { useState } from 'react';

import {
  AppAlert,
  AppBox,
  AppDialog,
  AppIconButton,
  AppStack,
  AppTypography
} from '@/shared/components';

import type { DashboardCreateRoomResult } from '../../types';

import {
  createRoomSuccessCopyButtonSx,
  createRoomSuccessCopyFieldSx,
  createRoomSuccessCopyValueSx,
  createRoomSuccessIconSx,
  createRoomSuccessLinkValueSx,
  createRoomSuccessTitleSx
} from './styles';

export interface CreateRoomSuccessDialogProps {
  readonly onClose: () => void;
  readonly onOpenRoom: () => void;
  readonly result: DashboardCreateRoomResult | null;
}

export const CreateRoomSuccessDialog = ({
  onClose,
  onOpenRoom,
  result
}: CreateRoomSuccessDialogProps) => {
  const [copiedValue, setCopiedValue] = useState<string | null>(null);

  if (!result) {
    return null;
  }

  const handleCopy = (value: string) => {
    navigator.clipboard
      .writeText(value)
      .then(() => {
        setCopiedValue(value);
      })
      .catch(() => {
        setCopiedValue(null);
      });
  };
  const handleCopyRoomCode = () => {
    handleCopy(result.roomCode);
  };
  const handleCopyInviteLink = () => {
    handleCopy(result.inviteLink);
  };

  return (
    <AppDialog
      aria-label="Room created successfully"
      cancelLabel="Close"
      confirmLabel="Open Room"
      maxWidth="sm"
      onCancel={onClose}
      onClose={onClose}
      onConfirm={onOpenRoom}
      open
    >
      <AppStack spacing={4}>
        <AppStack alignItems="center" spacing={2} textAlign="center">
          <AppBox sx={createRoomSuccessIconSx}>
            <CheckCircleRoundedIcon sx={{ fontSize: 34 }} />
          </AppBox>
          <AppStack spacing={1}>
            <AppTypography variant="h4">Room Created Successfully</AppTypography>
            <AppTypography color="text.secondary" variant="body2">
              Your technical estimation workspace is ready for your team.
            </AppTypography>
          </AppStack>
        </AppStack>
        <AppStack spacing={3}>
          <AppStack spacing={1}>
            <AppTypography color="text.secondary" variant="overline">
              Room Name
            </AppTypography>
            <AppTypography sx={createRoomSuccessTitleSx} title={result.roomName} variant="h6">
              {result.roomName}
            </AppTypography>
          </AppStack>
          <AppStack spacing={1}>
            <AppTypography color="text.secondary" variant="overline">
              Room Code
            </AppTypography>
            <AppBox alignItems="center" sx={createRoomSuccessCopyFieldSx}>
              <AppTypography
                color="primary.main"
                sx={createRoomSuccessCopyValueSx}
                title={result.roomCode}
                variant="h6"
              >
                {result.roomCode}
              </AppTypography>
              <AppIconButton
                aria-label="Copy room code"
                onClick={handleCopyRoomCode}
                size="small"
                sx={createRoomSuccessCopyButtonSx}
              >
                <ContentCopyRoundedIcon fontSize="small" />
              </AppIconButton>
            </AppBox>
            {copiedValue === result.roomCode ? (
              <AppTypography color="primary.main" variant="caption">
                Room code copied.
              </AppTypography>
            ) : null}
          </AppStack>
        </AppStack>
        <AppStack spacing={1.5}>
          <AppTypography color="text.secondary" variant="overline">
            Shareable Link
          </AppTypography>
          <AppBox sx={createRoomSuccessCopyFieldSx}>
            <AppTypography sx={createRoomSuccessLinkValueSx} title={result.inviteLink} variant="body2">
              {result.inviteLink}
            </AppTypography>
            <AppIconButton
              aria-label="Copy shareable link"
              onClick={handleCopyInviteLink}
              size="small"
              sx={createRoomSuccessCopyButtonSx}
            >
              <ContentCopyRoundedIcon fontSize="small" />
            </AppIconButton>
          </AppBox>
          {copiedValue === result.inviteLink ? (
            <AppTypography color="primary.main" variant="caption">
              Shareable link copied.
            </AppTypography>
          ) : null}
        </AppStack>
        {result.skippedRecipients.length > 0 ? (
          <AppAlert icon={<WarningAmberRoundedIcon fontSize="inherit" />} severity="warning">
            {result.skippedRecipients.length} participant invitation
            {result.skippedRecipients.length === 1 ? ' was' : 's were'} skipped.
          </AppAlert>
        ) : null}
      </AppStack>
    </AppDialog>
  );
};
