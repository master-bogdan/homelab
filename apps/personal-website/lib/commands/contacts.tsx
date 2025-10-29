import type React from "react"
import { Box, Typography, Link } from "@mui/material"

export const contactsCommand = (): React.ReactNode => {
  return (
    <Box>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00aaff", mb: 1 }}>
        Contact Information:
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff", mb: 0.5 }}>
        Email:{" "}
        <Link
          href="mailto:bogdan@example.com"
          sx={{ color: "#00ff00", textDecoration: "none", "&:hover": { textDecoration: "underline" } }}
        >
          bogdan@example.com
        </Link>
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff", mb: 0.5 }}>
        GitHub:{" "}
        <Link
          href="https://github.com/bogdan"
          target="_blank"
          rel="noopener noreferrer"
          sx={{ color: "#00ff00", textDecoration: "none", "&:hover": { textDecoration: "underline" } }}
        >
          github.com/bogdan
        </Link>
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff", mb: 0.5 }}>
        LinkedIn:{" "}
        <Link
          href="https://linkedin.com/in/bogdan"
          target="_blank"
          rel="noopener noreferrer"
          sx={{ color: "#00ff00", textDecoration: "none", "&:hover": { textDecoration: "underline" } }}
        >
          linkedin.com/in/bogdan
        </Link>
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff" }}>
        Twitter:{" "}
        <Link
          href="https://twitter.com/bogdan"
          target="_blank"
          rel="noopener noreferrer"
          sx={{ color: "#00ff00", textDecoration: "none", "&:hover": { textDecoration: "underline" } }}
        >
          @bogdan
        </Link>
      </Typography>
    </Box>
  )
}
