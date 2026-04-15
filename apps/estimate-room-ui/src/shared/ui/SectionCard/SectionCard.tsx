import type { PaperProps } from '@mui/material';
import type { ReactNode } from 'react';

import { AppStack } from '../AppStack';
import { AppTypography } from '../AppTypography';
import { SectionCardBody, SectionCardHeader, SectionCardRoot } from './styles';

export interface SectionCardProps extends Omit<PaperProps, 'title'> {
  readonly action?: ReactNode;
  readonly children?: ReactNode;
  readonly description?: string;
  readonly title: string;
}

export const SectionCard = ({
  action,
  children,
  description,
  title,
  ...paperProps
}: SectionCardProps) => (
  <SectionCardRoot elevation={0} {...paperProps}>
    <AppStack
      alignItems={{ sm: 'center' }}
      direction={{ sm: 'row' }}
      justifyContent="space-between"
      spacing={2}
    >
      <SectionCardHeader>
        <AppTypography component="h2" variant="h5">
          {title}
        </AppTypography>
        {description ? (
          <AppTypography color="text.secondary" variant="body2">
            {description}
          </AppTypography>
        ) : null}
      </SectionCardHeader>
      {action}
    </AppStack>
    {children ? <SectionCardBody>{children}</SectionCardBody> : null}
  </SectionCardRoot>
);
