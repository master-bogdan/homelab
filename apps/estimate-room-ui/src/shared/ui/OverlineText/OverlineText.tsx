import type { TypographyProps } from '@mui/material';

import { AppTypography } from '../AppTypography';

export type OverlineTextProps = Omit<TypographyProps, 'variant'>;

export const OverlineText = ({
  color = 'text.secondary',
  ...typographyProps
}: OverlineTextProps) => (
  <AppTypography color={color} variant="overline" {...typographyProps} />
);
