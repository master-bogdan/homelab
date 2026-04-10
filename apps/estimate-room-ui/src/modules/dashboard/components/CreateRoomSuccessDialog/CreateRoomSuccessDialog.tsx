import CheckCircleRoundedIcon from '@mui/icons-material/CheckCircleRounded';
import ContentCopyRoundedIcon from '@mui/icons-material/ContentCopyRounded';
import WarningAmberRoundedIcon from '@mui/icons-material/WarningAmberRounded';
import { Alert, Box, IconButton, Stack, Typography } from '@mui/material';
import { useState } from 'react';

import { AppDialog } from '@/shared/ui';

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

  const handleCopy = async (value: string) => {
    try {
      await navigator.clipboard.writeText(value);
      setCopiedValue(value);
    } catch {
      setCopiedValue(null);
    }
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
      <Stack spacing={4}>
        <Stack alignItems="center" spacing={2} textAlign="center">
          <Box sx={createRoomSuccessIconSx}>
            <CheckCircleRoundedIcon sx={{ fontSize: 34 }} />
          </Box>
          <Stack spacing={1}>
            <Typography variant="h4">Room Created Successfully</Typography>
            <Typography color="text.secondary" variant="body2">
              Your technical estimation workspace is ready for your team.
            </Typography>
          </Stack>
        </Stack>
        <Stack spacing={3}>
          <Stack spacing={1}>
            <Typography color="text.secondary" variant="overline">
              Room Name
            </Typography>
            <Typography sx={createRoomSuccessTitleSx} title={result.roomName} variant="h6">
              {result.roomName}
            </Typography>
          </Stack>
          <Stack spacing={1}>
            <Typography color="text.secondary" variant="overline">
              Room Code
            </Typography>
            <Box alignItems="center" sx={createRoomSuccessCopyFieldSx}>
              <Typography
                color="primary.main"
                sx={createRoomSuccessCopyValueSx}
                title={result.roomCode}
                variant="h6"
              >
                {result.roomCode}
              </Typography>
              <IconButton
                aria-label="Copy room code"
                onClick={() => {
                  void handleCopy(result.roomCode);
                }}
                size="small"
                sx={createRoomSuccessCopyButtonSx}
              >
                <ContentCopyRoundedIcon fontSize="small" />
              </IconButton>
            </Box>
            {copiedValue === result.roomCode ? (
              <Typography color="primary.main" variant="caption">
                Room code copied.
              </Typography>
            ) : null}
          </Stack>
        </Stack>
        <Stack spacing={1.5}>
          <Typography color="text.secondary" variant="overline">
            Shareable Link
          </Typography>
          <Box sx={createRoomSuccessCopyFieldSx}>
            <Typography sx={createRoomSuccessLinkValueSx} title={result.inviteLink} variant="body2">
              {result.inviteLink}
            </Typography>
            <IconButton
              aria-label="Copy shareable link"
              onClick={() => {
                void handleCopy(result.inviteLink);
              }}
              size="small"
              sx={createRoomSuccessCopyButtonSx}
            >
              <ContentCopyRoundedIcon fontSize="small" />
            </IconButton>
          </Box>
          {copiedValue === result.inviteLink ? (
            <Typography color="primary.main" variant="caption">
              Shareable link copied.
            </Typography>
          ) : null}
        </Stack>
        {result.skippedRecipients.length > 0 ? (
          <Alert icon={<WarningAmberRoundedIcon fontSize="inherit" />} severity="warning">
            {result.skippedRecipients.length} participant invitation
            {result.skippedRecipients.length === 1 ? ' was' : 's were'} skipped.
          </Alert>
        ) : null}
      </Stack>
    </AppDialog>
  );
};
