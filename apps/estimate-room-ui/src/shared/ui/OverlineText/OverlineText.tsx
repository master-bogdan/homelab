import { Typography } from '@mui/material';
import type { TypographyProps } from '@mui/material';

export type OverlineTextProps = Omit<TypographyProps, 'variant'>;

export const OverlineText = ({
  color = 'text.secondary',
  ...typographyProps
}: OverlineTextProps) => (
  <Typography color={color} variant="overline" {...typographyProps} />
);
