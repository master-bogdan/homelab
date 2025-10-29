import type React from "react"
import type { Metadata } from "next"
import { ThemeProvider } from "@mui/material/styles"
import CssBaseline from "@mui/material/CssBaseline"
import { theme } from "@/lib/theme"
import "@styles/globals.css" // Import globals.css here

export const metadata: Metadata = {
  title: "Bogdan Shchavinskyi - Terminal Portfolio",
  description: "Senior Software Engineer | Backend | Platform | DevOps",
  generator: "v0.app",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en">
      <body style={{ margin: 0, padding: 0 }}>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          {children}
        </ThemeProvider>
      </body>
    </html>
  )
}
