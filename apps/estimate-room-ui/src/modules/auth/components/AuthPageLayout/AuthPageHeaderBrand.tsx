import ArchitectureRoundedIcon from '@mui/icons-material/ArchitectureRounded';
import { Link as RouterLink } from 'react-router-dom';

import { AppRoutes } from '@/app/router/routePaths';
import { AppBox, AppLink, AppTypography } from '@/shared/components';

import {
  authPageLayoutHeaderRootSx,
  authPageLayoutHomeLinkSx
} from './styles';

export const AuthPageHeaderBrand = () => (
  <AppBox component="header" sx={authPageLayoutHeaderRootSx}>
    <AppLink
      color="inherit"
      component={RouterLink}
      sx={authPageLayoutHomeLinkSx}
      to={AppRoutes.LOGIN}
      underline="none"
    >
      <ArchitectureRoundedIcon color="primary" />
      <AppTypography color="text.primary" variant="h6">
        EstimateRoom
      </AppTypography>
    </AppLink>
  </AppBox>
);
