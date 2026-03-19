import type { PaperProps } from '@mui/material';
import { Stack, Typography } from '@mui/material';
import type { ReactNode } from 'react';

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
  <SectionCardRoot elevation={0} variant="outlined" {...paperProps}>
    <Stack
      alignItems={{ sm: 'center' }}
      direction={{ sm: 'row' }}
      justifyContent="space-between"
      spacing={2}
    >
      <SectionCardHeader>
        <Typography component="h2" variant="h5">
          {title}
        </Typography>
        {description ? (
          <Typography color="text.secondary" variant="body2">
            {description}
          </Typography>
        ) : null}
      </SectionCardHeader>
      {action}
    </Stack>
    {children ? <SectionCardBody>{children}</SectionCardBody> : null}
  </SectionCardRoot>
);
