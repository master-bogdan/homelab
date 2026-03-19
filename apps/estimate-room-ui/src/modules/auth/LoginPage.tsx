import CheckCircleOutlineRoundedIcon from '@mui/icons-material/CheckCircleOutlineRounded';
import { Chip, Stack, Typography } from '@mui/material';

import { SectionCard } from '@/shared/ui';

import { useLoginPage } from './hooks/useLoginPage';
import { authModuleNote } from './utils';

export const LoginPage = () => {
  const { readinessItems } = useLoginPage();

  return (
    <SectionCard
      description="Authentication is intentionally scaffolded without a fake sign-in implementation."
      sx={{ maxWidth: 720 }}
      title="Login"
    >
      <Stack direction="row" flexWrap="wrap" gap={1}>
        <Chip color="primary" label="Auth slice ready" variant="outlined" />
        <Chip color="secondary" label="Protected routes active" variant="outlined" />
        <Chip color="success" label="Backend integration pending" variant="outlined" />
      </Stack>
      <Stack component="ul" spacing={1.25} sx={{ m: 0, pl: 2 }}>
        {readinessItems.map((item) => (
          <Stack
            key={item}
            alignItems="flex-start"
            component="li"
            direction="row"
            spacing={1}
          >
            <CheckCircleOutlineRoundedIcon color="primary" fontSize="small" />
            <Typography variant="body2">{item}</Typography>
          </Stack>
        ))}
      </Stack>
      <Typography color="text.secondary" variant="body2">
        {authModuleNote}
      </Typography>
    </SectionCard>
  );
};
