import type React from "react"
import { Box, Typography } from "@mui/material"

export const projectsCommand = (): React.ReactNode => {
  return (
    <Box>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00aaff", mb: 1 }}>
        Featured Projects:
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff", ml: 2 }}>
        In development
      </Typography>
    </Box>
  )
}
