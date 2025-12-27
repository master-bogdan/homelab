import type React from "react"
import { Box, Typography, Link } from "@mui/material"

export const resumeCommand = (): React.ReactNode => {
  return (
    <Box>
      <Typography sx={{ fontFamily: "var(--font-mono)", fontSize: "14px", color: "#18c7ff", mb: 1 }}>Resume / CV</Typography>
      <Typography sx={{ fontFamily: "var(--font-mono)", fontSize: "14px", color: "#ffffff", mb: 1 }}>
        Download my resume:
      </Typography>
      <Link
        href="/Bogdan_Shchavinskyi_Senior_Software_Engineer_CV.pdf"
        download
        sx={{
          fontFamily: "var(--font-mono)",
          fontSize: "14px",
          color: "#00ff00",
          textDecoration: "none",
          "&:hover": { textDecoration: "underline" },
        }}
      >
        [Download PDF]
      </Link>
      <Typography sx={{ fontFamily: "var(--font-mono)", fontSize: "14px", color: "#888888", mt: 2 }}>
        Or view online at:
        <Link
          href="https://www.linkedin.com/in/b-shchavinskyi-fullstack"
          target="_blank"
          rel="noopener noreferrer"
          sx={{ color: "#00ff00", textDecoration: "none", "&:hover": { textDecoration: "underline" } }}
        >
          linkedin.com/in/b-shchavinskyi-fullstack
        </Link>
      </Typography>
    </Box>
  )
}
