"use client"

import type React from "react"
import { Box, Chip } from "@mui/material"

interface QuickCommandsProps {
  commands: string[]
  onCommand: (command: string) => void
}

export const QuickCommands: React.FC<QuickCommandsProps> = ({ commands, onCommand }) => {
  return (
    <Box
      sx={{
        display: "flex",
        flexWrap: "wrap",
        gap: 1,
        mb: 2,
        justifyContent: "flex-start",
      }}
    >
      {commands.map((cmd) => (
        <Chip
          key={cmd}
          label={cmd}
          onClick={() => onCommand(cmd)}
          sx={{
            fontFamily: "monospace",
            fontSize: "12px",
            backgroundColor: "#1a1a1a",
            color: "#00aaff",
            border: "1px solid #00aaff",
            cursor: "pointer",
            "&:hover": {
              backgroundColor: "#2a2a2a",
              borderColor: "#00ff00",
              color: "#00ff00",
            },
          }}
        />
      ))}
    </Box>
  )
}
