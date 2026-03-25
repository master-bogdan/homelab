import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';
import { Box, Stack, Typography } from '@mui/material';

export interface AuthIntroProps {
  readonly description: string;
  readonly title: string;
}

export const AuthIntro = ({ description, title }: AuthIntroProps) => (
  <Stack alignItems="center" spacing={2} sx={{ mb: 4.5, textAlign: 'center' }}>
    <Box
      sx={{
        alignItems: 'center',
        bgcolor: 'secondary.light',
        border: (theme) => `1px solid ${theme.app.borders.ghost}`,
        borderRadius: (theme) => Number(theme.shape.borderRadius) * 2,
        display: 'inline-flex',
        height: 52,
        justifyContent: 'center',
        width: 52
      }}
    >
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
