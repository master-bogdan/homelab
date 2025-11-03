"use client"

import type React from "react"
import { Box, Typography } from "@mui/material"
import { PROMPT_USER, PROMPT_HOST } from "@/lib/constants"

interface PromptProps {
  showCursor?: boolean
}

export const Prompt: React.FC<PromptProps> = ({ showCursor = true }) => {
  return (
    <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
      <Typography
        component="span"
        sx={{
          fontFamily: "monospace",
          fontSize: "14px",
          color: "#00ff00",
        }}
      >
        {PROMPT_USER}@{PROMPT_HOST}
      </Typography>
      <Typography
        component="span"
        sx={{
          fontFamily: "monospace",
          fontSize: "14px",
          color: "#ffffff",
        }}
      >
        :
      </Typography>
      <Typography
        component="span"
        sx={{
          fontFamily: "monospace",
          fontSize: "14px",
          color: "#00aaff",
        }}
      >
        ~
      </Typography>
      <Typography
        component="span"
        sx={{
          fontFamily: "monospace",
          fontSize: "14px",
          color: "#ffffff",
        }}
      >
        $
      </Typography>
      {showCursor && (
        <Box
          component="span"
          sx={{
            display: "inline-block",
            width: "10px",
            height: "18px",
            backgroundColor: "#00ff00",
            ml: 0.5,
            animation: "blink 1s step-end infinite",
            "@keyframes blink": {
              "0%, 50%": { opacity: 1 },
              "51%, 100%": { opacity: 0 },
            },
          }}
        />
      )}
    </Box>
  )
}
