import { useSettingsPage } from '@/modules/settings/hooks/useSettingsPage';
import { AppButton, AppStack, AppTypography, SectionCard } from '@/shared/components';

export const SettingsPage = () => {
  const { apiBaseUrl, environment, themeMode, toggleTheme, wsBaseUrl } = useSettingsPage();

  return (
    <SectionCard
      action={
        <AppButton onClick={toggleTheme} variant="contained">
          Toggle theme
        </AppButton>
      }
      description="Operational settings and environment-facing configuration can grow here."
      title="Settings"
    >
      <AppStack spacing={1}>
        <AppTypography color="text.secondary" variant="body2">
          Environment: {environment}
        </AppTypography>
        <AppTypography color="text.secondary" variant="body2">
          Active theme: {themeMode}
        </AppTypography>
        <AppTypography color="text.secondary" variant="body2">
          API base URL: {apiBaseUrl}
        </AppTypography>
        <AppTypography color="text.secondary" variant="body2">
          WebSocket URL: {wsBaseUrl}
        </AppTypography>
      </AppStack>
    </SectionCard>
  );
};
