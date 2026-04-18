import type { PropsWithChildren } from 'react';

import { AppBox } from '@/shared/components';

import {
  authPageLayoutInnerSx,
  authPageLayoutMainRootSx
} from './styles';

export const AuthPageContent = ({ children }: PropsWithChildren) => (
  <AppBox component="main" sx={authPageLayoutMainRootSx}>
    <AppBox sx={authPageLayoutInnerSx}>{children}</AppBox>
  </AppBox>
);
