import type React from "react"
import { Box, Typography, Link } from "@mui/material"

export const resumeCommand = (): React.ReactNode => {
  return (
    <Box>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00aaff", mb: 1 }}>Resume / CV</Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff", mb: 1 }}>
        Download my resume:
      </Typography>
      <Link
        href="/resume.pdf"
        download
        sx={{
          fontFamily: "monospace",
          fontSize: "14px",
          color: "#00ff00",
          textDecoration: "none",
          "&:hover": { textDecoration: "underline" },
        }}
      >
        [Download PDF]
      </Link>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#888888", mt: 2 }}>
        Or view online at: linkedin.com/in/bogdan
      </Typography>
    </Box>
  )
}
