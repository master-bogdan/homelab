import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';
import { Box, Stack, Typography } from '@mui/material';

import { authIntroIconSx, authIntroRootSx } from './styles';

export interface AuthIntroProps {
  readonly description: string;
  readonly title: string;
}

export const AuthIntro = ({ description, title }: AuthIntroProps) => (
  <Stack alignItems="center" spacing={2} sx={authIntroRootSx}>
    <Box sx={authIntroIconSx}>
      <ArchitectureRoundedIcon color="primary" />
    </Box>
    <Stack spacing={1}>
      <Typography component="h1" variant="h3">
        {title}
      </Typography>
      <Typography color="text.secondary" variant="body2">
        {description}
      </Typography>
    </Stack>
  </Stack>
);
