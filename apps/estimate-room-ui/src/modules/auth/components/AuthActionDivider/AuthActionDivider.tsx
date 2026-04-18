import { AppBox, AppStack, OverlineText } from '@/shared/components';

import { authActionDividerLineSx } from './styles';

interface AuthActionDividerProps {
  readonly label?: string;
}

export const AuthActionDivider = ({
  label = 'Or continue with'
}: AuthActionDividerProps) => (
  <AppStack alignItems="center" direction="row" spacing={2}>
    <AppBox sx={authActionDividerLineSx} />
    <OverlineText>{label}</OverlineText>
    <AppBox sx={authActionDividerLineSx} />
  </AppStack>
);
