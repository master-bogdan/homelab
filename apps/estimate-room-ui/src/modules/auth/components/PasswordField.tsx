import VisibilityOffRoundedIcon from '@mui/icons-material/VisibilityOffRounded';
import VisibilityRoundedIcon from '@mui/icons-material/VisibilityRounded';
import { IconButton, InputAdornment, TextField } from '@mui/material';
import type { TextFieldProps } from '@mui/material';
import { useState, type MouseEvent } from 'react';

export interface PasswordFieldProps extends Omit<TextFieldProps, 'type'> {
  readonly hideLabel?: string;
  readonly showLabel?: string;
}

export const PasswordField = ({
  hideLabel = 'Hide password',
  InputProps,
  showLabel = 'Show password',
  ...textFieldProps
}: PasswordFieldProps) => {
  const [isVisible, setIsVisible] = useState(false);

  const toggleVisibility = () => {
    setIsVisible((currentValue) => !currentValue);
  };

  const preventFocusLoss = (event: MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
  };

  return (
    <TextField
      {...textFieldProps}
      InputProps={{
        ...InputProps,
        endAdornment: (
          <>
            {InputProps?.endAdornment}
            <InputAdornment position="end">
              <IconButton
                aria-label={isVisible ? hideLabel : showLabel}
                aria-pressed={isVisible}
                edge="end"
                onClick={toggleVisibility}
                onMouseDown={preventFocusLoss}
                onMouseUp={preventFocusLoss}
              >
                {isVisible ? <VisibilityOffRoundedIcon /> : <VisibilityRoundedIcon />}
              </IconButton>
            </InputAdornment>
          </>
        )
      }}
      type={isVisible ? 'text' : 'password'}
    />
  );
};
