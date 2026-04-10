import CloseRoundedIcon from '@mui/icons-material/CloseRounded';
import { Dialog, IconButton, Typography } from '@mui/material';
import type { DialogProps } from '@mui/material';
import type { SxProps, Theme } from '@mui/material/styles';
import { useId } from 'react';
import type { ReactNode } from 'react';

import { AppButton } from '../AppButton';
import {
  AppDialogActionsRoot,
  AppDialogBodyRoot,
  AppDialogHeaderRoot,
  AppDialogTitleRoot
} from './styles';

export interface AppDialogProps extends Omit<DialogProps, 'title'> {
  readonly cancelDisabled?: boolean;
  readonly cancelLabel?: ReactNode;
  readonly children: ReactNode;
  readonly closeButtonAriaLabel?: string;
  readonly confirmDisabled?: boolean;
  readonly confirmLabel?: ReactNode;
  readonly confirmLoading?: boolean;
  readonly confirmLoadingText?: ReactNode;
  readonly confirmStartIcon?: ReactNode;
  readonly hideCancelButton?: boolean;
  readonly hideConfirmButton?: boolean;
  readonly onCancel?: () => void;
  readonly onConfirm?: () => void;
  readonly title?: ReactNode;
}

export const AppDialog = ({
  'aria-label': ariaLabel,
  cancelDisabled = false,
  cancelLabel = 'Cancel',
  children,
  closeButtonAriaLabel = 'Close dialog',
  confirmDisabled = false,
  confirmLabel = 'Apply',
  confirmLoading = false,
  confirmLoadingText,
  confirmStartIcon,
  fullWidth = true,
  hideCancelButton = false,
  hideConfirmButton = false,
  maxWidth = 'sm',
  onCancel,
  onClose,
  onConfirm,
  PaperProps,
  title,
  ...dialogProps
}: AppDialogProps) => {
  const titleId = useId();
  const hasTitle = Boolean(title);
  const basePaperSx: SxProps<Theme> = (theme) => ({
    backdropFilter: `blur(${theme.app.effects.backdropBlur})`,
    border: `1px solid ${theme.app.borders.ghost}`,
    borderRadius: theme.spacing(2),
    boxShadow: theme.app.effects.ambientShadow,
    overflow: 'hidden'
  });
  const inheritedPaperSx = PaperProps?.sx;
  const paperSx: SxProps<Theme> = Array.isArray(inheritedPaperSx)
    ? [basePaperSx, ...inheritedPaperSx]
    : inheritedPaperSx
      ? [basePaperSx, inheritedPaperSx]
      : [basePaperSx];

  const dialogAriaLabel =
    ariaLabel ?? (typeof title === 'string' || typeof title === 'number' ? String(title) : undefined);
  const paperAriaLabel = hasTitle ? PaperProps?.['aria-label'] : dialogAriaLabel;

  const handleCancel = () => {
    onCancel?.();
  };

  const handleClose = () => {
    onClose?.({}, 'escapeKeyDown');
  };

  return (
    <Dialog
      aria-labelledby={hasTitle ? titleId : undefined}
      PaperProps={{
        'aria-label': paperAriaLabel,
        elevation: 0,
        ...PaperProps,
        sx: paperSx
      }}
      fullWidth={fullWidth}
      maxWidth={maxWidth}
      onClose={onClose}
      {...dialogProps}
    >
      <AppDialogHeaderRoot>
        <AppDialogTitleRoot id={hasTitle ? titleId : undefined}>
          {typeof title === 'string' || typeof title === 'number' ? (
            <Typography component="span" fontWeight={800} variant="h5">
              {title}
            </Typography>
          ) : (
            title
          )}
        </AppDialogTitleRoot>
        <IconButton
          aria-label={closeButtonAriaLabel}
          disabled={cancelDisabled}
          onClick={handleClose}
          size="small"
        >
          <CloseRoundedIcon />
        </IconButton>
      </AppDialogHeaderRoot>

      <AppDialogBodyRoot>{children}</AppDialogBodyRoot>

      <AppDialogActionsRoot>
        {!hideCancelButton ? (
          <AppButton
            color="secondary"
            disabled={cancelDisabled}
            fullWidth
            onClick={handleCancel}
            variant="contained"
          >
            {cancelLabel}
          </AppButton>
        ) : null}
        {!hideConfirmButton ? (
          <AppButton
            disabled={confirmDisabled}
            fullWidth
            loading={confirmLoading}
            loadingText={confirmLoadingText}
            onClick={onConfirm}
            startIcon={confirmStartIcon}
            variant="contained"
          >
            {confirmLabel}
          </AppButton>
        ) : null}
      </AppDialogActionsRoot>
    </Dialog>
  );
};
