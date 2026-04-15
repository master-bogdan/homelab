import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';

import { AppBox, AppStack, AppTypography } from '@/shared/ui';

import { authIntroIconSx, authIntroRootSx } from './styles';

export interface AuthIntroProps {
  readonly description: string;
  readonly title: string;
}

export const AuthIntro = ({ description, title }: AuthIntroProps) => (
  <AppStack alignItems="center" spacing={2} sx={authIntroRootSx}>
    <AppBox sx={authIntroIconSx}>
      <ArchitectureRoundedIcon color="primary" />
    </AppBox>
    <AppStack spacing={1}>
      <AppTypography component="h1" variant="h3">
        {title}
      </AppTypography>
      <AppTypography color="text.secondary" variant="body2">
        {description}
      </AppTypography>
    </AppStack>
  </AppStack>
);
