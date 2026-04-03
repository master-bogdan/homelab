import { CircularProgress, Stack, Typography } from '@mui/material';
import type { StackProps, TypographyProps } from '@mui/material';
import type { ReactNode } from 'react';

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
  <Stack
    alignItems={alignItems}
    spacing={spacing}
    sx={{ textAlign, width: '100%', ...sx }}
    {...stackProps}
  >
    {isLoading ? (visual ?? <CircularProgress size={28} />) : visual}
    <Stack spacing={1.5} sx={{ maxWidth: 560 }}>
      <Typography component={titleComponent} variant={titleVariant}>
        {title}
      </Typography>
      {description ? (
        <Typography color="text.secondary" variant="body2">
          {description}
        </Typography>
      ) : null}
    </Stack>
    {action}
  </Stack>
);
