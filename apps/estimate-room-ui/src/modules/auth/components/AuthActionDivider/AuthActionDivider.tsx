import { Box, Stack } from '@mui/material';

import { OverlineText } from '@/shared/ui';

import { authActionDividerLineSx } from './styles';

export interface AuthActionDividerProps {
  readonly label?: string;
}

export const AuthActionDivider = ({
  label = 'Or continue with'
}: AuthActionDividerProps) => (
  <Stack alignItems="center" direction="row" spacing={2}>
    <Box sx={authActionDividerLineSx} />
    <OverlineText>{label}</OverlineText>
    <Box sx={authActionDividerLineSx} />
  </Stack>
);
