import { Box, DialogActions, DialogContent, styled } from '@mui/material';
import type { Theme } from '@mui/material/styles';

export const appDialogPaperSx = (theme: Theme) => ({
  backdropFilter: `blur(${theme.app.effects.backdropBlur})`,
  border: `1px solid ${theme.app.borders.ghost}`,
  borderRadius: theme.app.radii.xl,
  boxShadow: theme.app.effects.ambientShadow,
  overflow: 'hidden'
});

export const AppDialogHeaderRoot = styled(Box)(({ theme }) => ({
  alignItems: 'center',
  display: 'flex',
  gap: theme.spacing(2),
  justifyContent: 'space-between',
  minHeight: 88,
  paddingBlock: theme.spacing(2.5),
  paddingInline: theme.spacing(4)
}));

export const AppDialogTitleRoot = styled(Box)({
  flex: 1,
  minWidth: 0
});

export const AppDialogBodyRoot = styled(DialogContent)(({ theme }) => ({
  overflow: 'visible',
  paddingBottom: theme.spacing(4),
  paddingInline: theme.spacing(4),
  paddingTop: theme.spacing(4)
}));

export const AppDialogActionsRoot = styled(DialogActions)(({ theme }) => ({
  display: 'grid',
  gap: theme.spacing(2),
  gridAutoColumns: '1fr',
  gridAutoFlow: 'column',
  paddingBottom: theme.spacing(3),
  paddingInline: theme.spacing(4),
  paddingTop: theme.spacing(3)
}));
