import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';

import { AppBox, AppStack, AppTypography } from '@/shared/ui';

import { authPageHeaderIconSx, authPageHeaderRootSx } from './styles';

interface AuthPageHeaderProps {
  readonly description: string;
  readonly title: string;
}

export const AuthPageHeader = ({ description, title }: AuthPageHeaderProps) => (
  <AppStack alignItems="center" spacing={2} sx={authPageHeaderRootSx}>
    <AppBox sx={authPageHeaderIconSx}>
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
