import type React from "react"
import { Box, Typography } from "@mui/material"
import { commandRegistry } from "../commandRegistry"
import { COMMANDS } from '../constants'

export const helpCommand = (args: string[]): React.ReactNode => {
  const showHidden = args.includes("--hidden")

  if (showHidden) {
    const hiddenCommands = Object.values(commandRegistry).filter((cmd) => cmd.hidden)

    return (
      <Box>
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ff00ff", mb: 1 }}>
          🎉 Secret Hidden Commands:
        </Typography>
        {hiddenCommands.map((cmd) => (
          <Box key={cmd.name} sx={{ display: "flex", mb: 0.5 }}>
            <Typography
              sx={{
                fontFamily: "monospace",
                fontSize: "14px",
                color: "#ff00ff",
                minWidth: "120px",
              }}
            >
              {cmd.name}
            </Typography>
            <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff" }}>
              - {cmd.description}
            </Typography>
          </Box>
        ))}
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#888888", mt: 2 }}>
          Congratulations on finding the secret! 🎊
        </Typography>
      </Box>
    )
  }

  return (
    <Box>
      <Typography
        key="available-commands-title"
        sx={{
          fontFamily: "monospace",
          fontSize: "14px",
          color: "#00aaff",
          mb: 1,
        }}>
        Available Commands:
      </Typography>
      {COMMANDS.map((cmd) => (
        <Box key={cmd.name} sx={{ display: "flex", mb: 0.5 }}>
          <Typography
            sx={{
              fontFamily: "monospace",
              fontSize: "14px",
              color: "#00ff00",
              minWidth: "120px",
            }}
          >
            {cmd.name}
          </Typography>
          <Typography
            sx={{
              fontFamily: "monospace",
              fontSize: "14px",
              color: "#ffffff"
            }}>
            - {cmd.description}
          </Typography>
        </Box>
      ))}
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#888888", mt: 2 }}>
        Hint: Try exploring for hidden commands...
      </Typography>
    </Box>
  )
}
