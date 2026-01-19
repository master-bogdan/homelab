"use client"

import { createTheme } from "@mui/material/styles"

export const theme = createTheme({
  palette: {
    mode: "dark",
    primary: {
      main: "#00ff00",
    },
    secondary: {
      main: "#18c7ff",
    },
    background: {
      default: "#0a0a0a",
      paper: "#1a1a1a",
    },
    text: {
      primary: "#00ff00",
      secondary: "#18c7ff",
    },
  },
  typography: {
    fontFamily: "var(--font-mono)",
    fontSize: 14,
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          backgroundColor: "#0a0a0a",
          color: "#00ff00",
          fontFamily: "var(--font-mono)",
        },
      },
    },
  },
})
