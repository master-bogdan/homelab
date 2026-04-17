import { Link as RouterLink } from 'react-router-dom';
import type { To } from 'react-router-dom';

import { AppLink } from '@/shared/ui';

import { AuthInlinePrompt } from '../AuthInlinePrompt';

interface AuthPageFooterProps {
  readonly linkLabel: string;
  readonly prompt: string;
  readonly to: To;
}

export const AuthPageFooter = ({
  linkLabel,
  prompt,
  to
}: AuthPageFooterProps) => (
  <AuthInlinePrompt>
    {prompt}{' '}
    <AppLink color="primary" component={RouterLink} to={to} underline="none">
      {linkLabel}
    </AppLink>
  </AuthInlinePrompt>
);
