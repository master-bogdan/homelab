import type { SxProps, Theme } from '@mui/material/styles';

export const createRoomSuccessIconSx: SxProps<Theme> = {
  alignItems: 'center',
  bgcolor: (theme) => `${theme.palette.primary.main}1f`,
  borderRadius: 2.5,
  color: 'primary.main',
  display: 'flex',
  height: 68,
  justifyContent: 'center',
  width: 68
};

export const createRoomSuccessTitleSx: SxProps<Theme> = {
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap'
};

export const createRoomSuccessCopyFieldSx: SxProps<Theme> = {
  alignItems: 'center',
  bgcolor: 'secondary.light',
  borderRadius: 2,
  display: 'flex',
  gap: 1,
  minWidth: 0,
  px: 2,
  py: 1.25
};

export const createRoomSuccessCopyValueSx: SxProps<Theme> = {
  flex: 1,
  fontWeight: 700,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap'
};

export const createRoomSuccessLinkValueSx: SxProps<Theme> = {
  flex: 1,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap'
};

export const createRoomSuccessCopyButtonSx: SxProps<Theme> = {
  color: 'primary.main',
  flexShrink: 0
};
