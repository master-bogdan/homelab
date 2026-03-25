import { Box, Stack } from '@mui/material';

import { OverlineText } from '@/shared/ui';

export interface AuthActionDividerProps {
  readonly label?: string;
}

const dividerLineSx = {
  bgcolor: (theme: { app: { borders: { ghost: string } } }) => theme.app.borders.ghost,
  flex: 1,
  height: 1
} as const;

export const AuthActionDivider = ({
  label = 'Or continue with'
}: AuthActionDividerProps) => (
  <Stack alignItems="center" direction="row" spacing={2}>
    <Box sx={dividerLineSx} />
    <OverlineText>{label}</OverlineText>
    <Box sx={dividerLineSx} />
  </Stack>
);
