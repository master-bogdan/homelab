import type { PropsWithChildren } from 'react';

import { AppBox } from '@/shared/components';

import {
  type AuthPageLayoutPattern,
  authPageLayoutRootSx
} from './styles';
import { AuthPageBackdrop } from './AuthPageBackdrop';
import { AuthPageContent } from './AuthPageContent';
import { AuthPageFooter } from './AuthPageFooter';
import { AuthPageHeaderBrand } from './AuthPageHeaderBrand';

interface AuthPageLayoutProps extends PropsWithChildren {
  readonly pattern?: AuthPageLayoutPattern;
}

export const AuthPageLayout = ({
  children,
  pattern = 'ambient'
}: AuthPageLayoutProps) => (
  <AppBox sx={authPageLayoutRootSx}>
    <AuthPageBackdrop pattern={pattern} />
    <AuthPageHeaderBrand />
    <AuthPageContent>{children}</AuthPageContent>
    <AuthPageFooter />
  </AppBox>
);
