import CheckCircleRoundedIcon from '@mui/icons-material/CheckCircleRounded';
import RadioButtonUncheckedRoundedIcon from '@mui/icons-material/RadioButtonUncheckedRounded';
import { Box, Stack, Typography } from '@mui/material';

import { OverlineText } from '@/shared/ui';

const getRecommendations = (password: string) => [
  {
    isMet: password.length >= 8,
    label: '8+ characters'
  },
  {
    isMet: /[0-9]/u.test(password),
    label: 'One number'
  },
  {
    isMet: /[A-Z]/u.test(password),
    label: 'Uppercase letter'
  },
  {
    isMet: /[^A-Za-z0-9]/u.test(password),
    label: 'Special symbol'
  }
];

export interface PasswordRecommendationsProps {
  readonly password: string;
}

export const PasswordRecommendations = ({
  password
}: PasswordRecommendationsProps) => {
  const recommendations = getRecommendations(password);

  return (
    <Box
      sx={{
        bgcolor: 'secondary.light',
        borderRadius: 0.5,
        px: 2,
        py: 1.75
      }}
    >
      <OverlineText sx={{ mb: 1.5 }}>Recommendations</OverlineText>
      <Box
        sx={{
          columnGap: 3,
          display: 'grid',
          gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, minmax(0, 1fr))' },
          rowGap: 1
        }}
      >
        {recommendations.map((recommendation) => {
          const Icon = recommendation.isMet
            ? CheckCircleRoundedIcon
            : RadioButtonUncheckedRoundedIcon;

          return (
            <Stack
              key={recommendation.label}
              alignItems="center"
              direction="row"
              spacing={1}
            >
              <Icon
                color={recommendation.isMet ? 'primary' : 'disabled'}
                fontSize="small"
              />
              <Typography color="text.secondary" variant="caption">
                {recommendation.label}
              </Typography>
            </Stack>
          );
        })}
      </Box>
    </Box>
  );
};
