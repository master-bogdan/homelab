import type React from "react"
import { Box, Typography } from "@mui/material"

export const whoamiCommand = (): React.ReactNode => {
  return (
    <Box>
      <Typography sx={{ fontFamily: "var(--font-mono)", fontSize: "14px", color: "#00ff00", mb: 1 }}>
        Bogdan Shchavinskyi
      </Typography>
      <Typography sx={{ fontFamily: "var(--font-mono)", fontSize: "14px", color: "#ffffff", mb: 0.5 }}>
        Senior Software Engineer
      </Typography>
      <Typography sx={{ fontFamily: "var(--font-mono)", fontSize: "14px", color: "#ffffff", mb: 0.5 }}>
        Specializations: Backend | Platform | DevOps
      </Typography>
      <Typography sx={{ fontFamily: "var(--font-mono)", fontSize: "14px", color: "#ffffff", mb: 0.5 }}>
        Expert in: Node.js, TypeScript, Go (Golang)
      </Typography>
      <Typography sx={{ fontFamily: "var(--font-mono)", fontSize: "14px", color: "#ffffff" }}>
        Cloud: AWS, Kubernetes
      </Typography>
    </Box>
  )
}
