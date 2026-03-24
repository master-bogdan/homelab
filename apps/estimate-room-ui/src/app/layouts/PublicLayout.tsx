import { Box, Container, Stack, Typography } from '@mui/material';
import { Outlet } from 'react-router-dom';

import { appConfig } from '@/shared/config/env';

import styles from './PublicLayout.module.scss';

export const PublicLayout = () => (
  <Box className={styles.root}>
    <Container maxWidth="lg" sx={{ minHeight: '100vh', py: { xs: 6, md: 10 } }}>
      <Stack spacing={4}>
        <Stack spacing={1}>
          <Typography component="p" color="primary.main" variant="overline">
            Estimate Platform
          </Typography>
          <Typography component="h1" variant="h3">
            {appConfig.appName}
          </Typography>
          <Typography color="text.secondary" maxWidth={720} variant="body1">
            Production-focused React scaffold for room estimates, history review,
            team workflows, and backend API integration.
          </Typography>
        </Stack>
        <Outlet />
      </Stack>
    </Container>
  </Box>
);
