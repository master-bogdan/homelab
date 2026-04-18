import { Box, Paper, Stack, styled } from '@mui/material';

export const SectionCardRoot = styled(Paper)(({ theme }) => ({
  borderRadius: theme.app.radii.lg,
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
