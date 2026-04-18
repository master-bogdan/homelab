import type { ReactNode } from 'react';

import { AppTypography } from '@/shared/components';

import { authInlinePromptSx } from './styles';

interface AuthInlinePromptProps {
  readonly children: ReactNode;
}

export const AuthInlinePrompt = ({ children }: AuthInlinePromptProps) => (
  <AppTypography sx={authInlinePromptSx} variant="body2">
    {children}
  </AppTypography>
);
