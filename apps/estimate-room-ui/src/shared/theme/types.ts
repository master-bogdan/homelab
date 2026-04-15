export type ThemeMode = 'dark' | 'light';

export interface AppThemeSurfaces {
  readonly base: string;
  readonly section: string;
  readonly card: string;
  readonly cardHover: string;
  readonly inset: string;
  readonly overlay: string;
  readonly bright: string;
  readonly well: string;
  readonly rowAlternate: string;
}

export interface AppThemeBorders {
  readonly ghost: string;
  readonly focusRing: string;
}

export interface AppThemeGradients {
  readonly primary: string;
}

export interface AppThemeEffects {
  readonly ambientShadow: string;
  readonly backdropBlur: string;
}

export interface AppThemeBackgrounds {
  readonly authAmbient: string;
  readonly authDots: string;
  readonly body: string;
  readonly drawer: string;
}

export interface AppThemeRadii {
  readonly sm: number;
  readonly md: number;
  readonly lg: number;
  readonly xl: number;
  readonly pill: number;
  readonly circle: string;
}

export interface AppThemeStateLayers {
  readonly primarySoft: string;
  readonly secondaryPanel: string;
}

export interface AppThemeLayout {
  readonly sectionGap: number;
  readonly pageGap: number;
}

export interface AppThemeTokens {
  readonly backgrounds: AppThemeBackgrounds;
  readonly surfaces: AppThemeSurfaces;
  readonly borders: AppThemeBorders;
  readonly gradients: AppThemeGradients;
  readonly effects: AppThemeEffects;
  readonly radii: AppThemeRadii;
  readonly stateLayers: AppThemeStateLayers;
  readonly layout: AppThemeLayout;
}

declare module '@mui/material/styles' {
  interface Theme {
    readonly app: AppThemeTokens;
  }

  interface ThemeOptions {
    readonly app?: AppThemeTokens;
  }
}
