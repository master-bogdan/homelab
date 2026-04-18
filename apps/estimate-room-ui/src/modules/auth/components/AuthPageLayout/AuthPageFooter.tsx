import { AppBox, AppLink, AppStack, OverlineText } from '@/shared/components';

import {
  authPageLayoutFooterRootSx,
  authPageLayoutFooterStackSx,
  authPageLayoutUtilityLinkSx
} from './styles';

const footerLinks = ['Privacy', 'Terms', 'Security'] as const;

export const AuthPageFooter = () => {
  const currentYear = new Date().getFullYear();

  return (
    <AppBox component="footer" sx={authPageLayoutFooterRootSx}>
      <AppStack
        alignItems={{ xs: 'flex-start', md: 'center' }}
        direction={{ xs: 'column', md: 'row' }}
        justifyContent="space-between"
        spacing={2}
        sx={authPageLayoutFooterStackSx}
      >
        <OverlineText>© {currentYear} EstimateRoom. All rights reserved.</OverlineText>
        <AppStack direction="row" spacing={3}>
          {footerLinks.map((label) => (
            <AppLink
              key={label}
              color="text.secondary"
              href="#"
              onClick={(event) => event.preventDefault()}
              sx={authPageLayoutUtilityLinkSx}
              underline="always"
              variant="overline"
            >
              {label}
            </AppLink>
          ))}
        </AppStack>
      </AppStack>
    </AppBox>
  );
};
