import { Box, Paper, Stack, styled } from '@mui/material';

export const SectionCardRoot = styled(Paper)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(2),
  height: '100%',
  padding: theme.spacing(3),
  borderRadius: `calc(${theme.shape.borderRadius} * 1.25)`,
  backgroundImage: 'none'
}));

export const SectionCardHeader = styled(Stack)(({ theme }) => ({
  gap: theme.spacing(0.75)
}));

export const SectionCardBody = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  gap: 'inherit'
});
