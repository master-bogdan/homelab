import type React from "react"
import { Box, Typography, Link } from "@mui/material"

export const contactsCommand = (): React.ReactNode => {
  return (
    <Box>
      <Typography sx={{ fontFamily: "var(--font-mono)", fontSize: "14px", color: "#00aaff", mb: 1 }}>
        Contact Information:
      </Typography>
      <Typography sx={{ fontFamily: "var(--font-mono)", fontSize: "14px", color: "#ffffff", mb: 0.5 }}>
        Email:{" "}
        <Link
          href="mailto:bshchavinskyi@gmail.com"
          sx={{ color: "#00ff00", textDecoration: "none", "&:hover": { textDecoration: "underline" } }}
        >
          bshchavinskyi@gmail.com
        </Link>
      </Typography>
      <Typography sx={{ fontFamily: "var(--font-mono)", fontSize: "14px", color: "#ffffff", mb: 0.5 }}>
        GitHub:{" "}
        <Link
          href="https://github.com/master-bogdan"
          target="_blank"
          rel="noopener noreferrer"
          sx={{ color: "#00ff00", textDecoration: "none", "&:hover": { textDecoration: "underline" } }}
        >
          github.com/master-bogdan
        </Link>
      </Typography>
      <Typography sx={{ fontFamily: "var(--font-mono)", fontSize: "14px", color: "#ffffff", mb: 0.5 }}>
        LinkedIn:{" "}
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
