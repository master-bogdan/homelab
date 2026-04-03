import type { PropsWithChildren } from 'react';
import { Box, Container, Toolbar } from '@mui/material';

export interface ContentShellProps extends PropsWithChildren {
  readonly drawerWidth: number;
}

export const ContentShell = ({ children, drawerWidth }: ContentShellProps) => (
  <Box
    component="main"
    sx={{
      backgroundColor: 'background.default',
      flexGrow: 1,
      minWidth: 0,
      width: { lg: `calc(100% - ${drawerWidth}px)` }
    }}
  >
    <Toolbar />
    <Container
      maxWidth="xl"
      sx={{
        py: { xs: 3, md: 5 },
        px: { xs: 2.5, md: 4 }
      }}
    >
      {children}
    </Container>
  </Box>
);
