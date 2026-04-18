import type { SxProps, Theme } from '@mui/material/styles';

export const authBackToSignInLinkSx: SxProps<Theme> = {
  alignItems: 'center',
  display: 'inline-flex',
  gap: 1
};

export const authBackToSignInCenteredLinkSx: SxProps<Theme> = {
  ...authBackToSignInLinkSx,
  justifyContent: 'center'
};

export const authBackToSignInFormLinkSx: SxProps<Theme> = {
  ...authBackToSignInLinkSx,
  mx: 'auto'
};
