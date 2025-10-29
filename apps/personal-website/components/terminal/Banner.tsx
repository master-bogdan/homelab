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
          fontSize: "12px",
          color: "#00aaff",
          lineHeight: 1.2,
          whiteSpace: "pre",
          margin: 0,
        }}
      >
        {ASCII_BANNER}
      </Typography>
    </Box>
  )
}
