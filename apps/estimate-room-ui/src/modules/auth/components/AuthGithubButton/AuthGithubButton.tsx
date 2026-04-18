import GitHubIcon from '@mui/icons-material/GitHub';

import { AppButton } from '@/shared/components';

interface AuthGithubButtonProps {
  readonly children: string;
  readonly disabled: boolean;
  readonly loading: boolean;
  readonly loadingText?: string;
  readonly onClick: () => void;
}

export const AuthGithubButton = ({
  children,
  disabled,
  loading,
  loadingText = 'Redirecting to GitHub...',
  onClick
}: AuthGithubButtonProps) => (
  <AppButton
    color="secondary"
    disabled={disabled}
    fullWidth
    loading={loading}
    loadingText={loadingText}
    onClick={onClick}
    startIcon={<GitHubIcon />}
    type="button"
    variant="contained"
  >
    {children}
  </AppButton>
);
