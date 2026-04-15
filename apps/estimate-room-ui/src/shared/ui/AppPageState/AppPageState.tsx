import type { StackProps, TypographyProps } from '@mui/material';
import type { ReactNode } from 'react';

import { AppProgress } from '../AppProgress';
import { AppStack } from '../AppStack';
import { AppTypography } from '../AppTypography';

export interface AppPageStateProps extends Omit<StackProps, 'title'> {
  readonly action?: ReactNode;
  readonly description?: ReactNode;
  readonly isLoading?: boolean;
  readonly title: ReactNode;
  readonly titleComponent?: TypographyProps['component'];
  readonly titleVariant?: TypographyProps['variant'];
  readonly visual?: ReactNode;
}

export const AppPageState = ({
  action,
  alignItems = 'center',
  description,
  isLoading = false,
  spacing = 3,
  sx,
  textAlign = 'center',
  title,
  titleComponent = 'h2',
  titleVariant = 'h5',
  visual,
  ...stackProps
}: AppPageStateProps) => (
  <AppStack
    alignItems={alignItems}
    spacing={spacing}
    sx={{ textAlign, width: '100%', ...sx }}
    {...stackProps}
  >
    {isLoading ? (visual ?? <AppProgress size={28} />) : visual}
    <AppStack spacing={1.5} sx={{ maxWidth: 560 }}>
      <AppTypography component={titleComponent} variant={titleVariant}>
        {title}
      </AppTypography>
      {description ? (
        <AppTypography color="text.secondary" variant="body2">
          {description}
        </AppTypography>
      ) : null}
    </AppStack>
    {action}
  </AppStack>
);
