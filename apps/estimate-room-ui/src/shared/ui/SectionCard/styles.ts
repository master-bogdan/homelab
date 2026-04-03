import { Box, Paper, Stack, styled } from '@mui/material';

export const SectionCardRoot = styled(Paper)(({ theme }) => ({
  // MUI types borderRadius as number | string, so cast before applying arithmetic.
  borderRadius:
    theme.palette.mode === 'light'
      ? Number(theme.shape.borderRadius) * 1.25
      : Number(theme.shape.borderRadius) * 2,
  backgroundColor: theme.app.surfaces.card,
  border:
    theme.palette.mode === 'light' ? `1px solid ${theme.app.borders.ghost}` : 'none',
  boxShadow: 'none',
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(2),
  height: '100%',
  padding: theme.spacing(3),
  transition: theme.transitions.create(['background-color', 'border-color'])
}));

export const SectionCardHeader = styled(Stack)(({ theme }) => ({
  gap: theme.spacing(0.75),
  minWidth: 0
}));

export const SectionCardBody = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  gap: 'inherit'
});
