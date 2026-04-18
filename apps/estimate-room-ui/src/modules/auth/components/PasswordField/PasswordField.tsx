import VisibilityOffRoundedIcon from '@mui/icons-material/VisibilityOffRounded';
import VisibilityRoundedIcon from '@mui/icons-material/VisibilityRounded';
import { useState, type MouseEvent } from 'react';

import { AppIconButton, AppInputAdornment, AppTextField } from '@/shared/components';
import type { AppTextFieldProps } from '@/shared/components';

import { passwordFieldToggleButtonSx } from './styles';

interface PasswordFieldProps extends Omit<AppTextFieldProps, 'type'> {
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
    <AppTextField
      {...textFieldProps}
      InputProps={{
        ...InputProps,
        endAdornment: (
          <>
            {InputProps?.endAdornment}
            <AppInputAdornment position="end">
              <AppIconButton
                aria-label={isVisible ? hideLabel : showLabel}
                aria-pressed={isVisible}
                edge="end"
                onClick={toggleVisibility}
                onMouseDown={preventFocusLoss}
                onMouseUp={preventFocusLoss}
                sx={passwordFieldToggleButtonSx}
              >
                {isVisible ? <VisibilityOffRoundedIcon /> : <VisibilityRoundedIcon />}
              </AppIconButton>
            </AppInputAdornment>
          </>
        )
      }}
      type={isVisible ? 'text' : 'password'}
    />
  );
};
