import { Button, Stack, Typography } from '@mui/material';

import { SectionCard } from '@/shared/ui';

import { useSettingsPage } from './hooks/useSettingsPage';

export const SettingsPage = () => {
  const { apiBaseUrl, environment, themeMode, toggleTheme, wsBaseUrl } = useSettingsPage();

  return (
    <SectionCard
      action={
        <Button onClick={toggleTheme} variant="contained">
          Toggle theme
        </Button>
      }
      description="Operational settings and environment-facing configuration can grow here."
      title="Settings"
    >
      <Stack spacing={1}>
        <Typography color="text.secondary" variant="body2">
          Environment: {environment}
        </Typography>
        <Typography color="text.secondary" variant="body2">
          Active theme: {themeMode}
        </Typography>
        <Typography color="text.secondary" variant="body2">
          API base URL: {apiBaseUrl}
        </Typography>
        <Typography color="text.secondary" variant="body2">
          WebSocket URL: {wsBaseUrl}
        </Typography>
      </Stack>
    </SectionCard>
  );
};
