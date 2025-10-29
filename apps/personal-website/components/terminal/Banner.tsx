"use client"

import type React from "react"
import { Box, Typography } from "@mui/material"
import { ASCII_BANNER } from "@/lib/constants"

export const Banner: React.FC = () => {
  return (
    <Box sx={{ mb: 2 }}>
      <Typography
        component="pre"
        sx={{
          fontFamily: "monospace",
          fontSize: { xs: "6px", sm: "8px", md: "12px" },
          color: "#00aaff",
          lineHeight: 1.2,
          whiteSpace: "pre",
          margin: 0,
          overflow: "hidden",
        }}
      >
        {ASCII_BANNER}
      </Typography>
    </Box>
  )
}
