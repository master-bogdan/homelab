import { TextField } from '@mui/material';
import type { TextFieldProps } from '@mui/material';

export type AppTextFieldProps = TextFieldProps & {
  readonly reserveHelperTextSpace?: boolean;
};

export const AppTextField = ({
  fullWidth = true,
  helperText,
  reserveHelperTextSpace = true,
  variant = 'outlined',
  ...textFieldProps
}: AppTextFieldProps) => (
  <TextField
    fullWidth={fullWidth}
    helperText={helperText ?? (reserveHelperTextSpace ? ' ' : undefined)}
    variant={variant}
    {...textFieldProps}
  />
);
