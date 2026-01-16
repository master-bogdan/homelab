import type React from "react"
import { Box } from "@mui/material"

export const Banner: React.FC = () => {
  return (
    <Box
      component="img"
      src="banner.svg"
      alt="Bogdan Shchavinskyi banner"
      sx={{
        display: "block",
        width: "100%",
        maxWidth: "860px",
        height: "auto",
        mb: 2,
      }}
    />
  )
}
