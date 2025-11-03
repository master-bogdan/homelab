import type React from "react"
import { Box, Typography } from "@mui/material"

export const neofetchCommand = (): React.ReactNode => {
   const uptimeSec = typeof performance !== "undefined" && typeof performance.now === "function" ? Math.floor(performance.now() / 1000) : 0

  return (
    <Box sx={{ display: "flex", gap: 2, flexWrap: "wrap" }}>
      <Box sx={{ flex: "0 0 auto" }}>
        <Typography component="pre" sx={{ fontFamily: "monospace", fontSize: "14px", color: "#dd4814", m: 0 }}>
          {`       _,met$$$$$gg.
    ,g$$$$$$$$$$$$$$$P.
  ,g$$P"     """Y$$.".
 ,$$P'              \`$$$.
',$$P       ,ggs.     \`$$b:
\`d$$'     ,$P"'   .    $$$
 $$P      d$'     ,    $$P
 $$:      $$.   -    ,d$$'
 $$;      Y$b._   _,d$P'
 Y$$.    \`.\`"Y$$$$P"'
 \`$$b      "-.__
  \`Y$$
   \`Y$$.
     \`$$b.
       \`Y$$b.
          \`"Y$b._
              \`""""`}
        </Typography>
      </Box>
      <Box sx={{ flex: "1 1 auto", minWidth: "300px" }}>
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00aaff", fontWeight: "bold" }}>
          bogdan@homelab
        </Typography>
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#888888" }}>---------------</Typography>
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff" }}>
          <Box component="span" sx={{ color: "#00aaff" }}>
            OS:
          </Box>{" "}
          Ubuntu 22.04 LTS x86_64
        </Typography>
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff" }}>
          <Box component="span" sx={{ color: "#00aaff" }}>
            Host:
          </Box>{" "}
          Personal Portfolio
        </Typography>
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff" }}>
          <Box component="span" sx={{ color: "#00aaff" }}>
            Kernel:
          </Box>{" "}
          6.2.0-terminal
        </Typography>
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff" }}>
          <Box component="span" sx={{ color: "#00aaff" }}>
            Uptime:
          </Box>{" "}
          {uptimeSec} seconds
        </Typography>
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff" }}>
          <Box component="span" sx={{ color: "#00aaff" }}>
            Shell:
          </Box>{" "}
          bash 5.1.16
        </Typography>
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff" }}>
          <Box component="span" sx={{ color: "#00aaff" }}>
            Terminal:
          </Box>{" "}
          xterm-256color
        </Typography>
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff" }}>
          <Box component="span" sx={{ color: "#00aaff" }}>
            CPU:
          </Box>{" "}
          TypeScript Engine (8) @ 3.9GHz
        </Typography>
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff" }}>
          <Box component="span" sx={{ color: "#00aaff" }}>
            Memory:
          </Box>{" "}
          {Math.floor(Math.random() * 2000 + 1000)}MiB / 16384MiB
        </Typography>
      </Box>
    </Box>
  )
}

export const sudoCommand = (): React.ReactNode => {
  return (
    <Box>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ff0000", mb: 1 }}>
        [sudo] password for bogdan:
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff", mb: 1 }}>
        Sorry, try again.
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#888888" }}>
        sudo: 3 incorrect password attempts
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00ff00", mt: 1 }}>
        Just kidding! You don&apos;t need sudo here. 😉
      </Typography>
    </Box>
  )
}

export const hackCommand = (): React.ReactNode => {
  const hackLines = [
    "Initializing hack sequence...",
    "Connecting to mainframe...",
    "Bypassing firewall...",
    "Decrypting passwords...",
    "Accessing database...",
    "Downloading files...",
    "Covering tracks...",
    "",
    "HACK COMPLETE! 🎉",
    "",
    "Just kidding. Please don't hack anyone. 🙃",
  ]

  return (
    <Box>
      {hackLines.map((line, index) => (
        <Typography
          key={index}
          sx={{
            fontFamily: "monospace",
            fontSize: "14px",
            color: index === hackLines.length - 3 ? "#00ff00" : "#00aaff",
            mb: 0.5,
          }}
        >
          {line}
        </Typography>
      ))}
    </Box>
  )
}

export const fortuneCommand = (): React.ReactNode => {
  const fortunes = [
    "You will write bug-free code today... just kidding!",
    "A merge conflict approaches. Prepare yourself.",
    "Your next commit will be legendary.",
    "The production server is stable... for now.",
    "Coffee levels are optimal for coding.",
    "A wild segfault appears!",
    "Your code will compile on the first try. (Unlikely)",
    "The cloud is just someone else's computer.",
  ]
  const randomFortune = fortunes[Math.floor(Math.random() * fortunes.length)]

  return (
    <Box>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00ff00" }}>{randomFortune}</Typography>
    </Box>
  )
}

export const cowsayCommand = (args: string[]): React.ReactNode => {
  const message = args.length > 0 ? args.join(" ") : "Hello from the terminal!"
  const messageLength = message.length

  return (
    <Box>
      <Typography component="pre" sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff" }}>
        {`
 ${"_".repeat(messageLength + 2)}
< ${message} >
 ${"-".repeat(messageLength + 2)}
        \\   ^__^
         \\  (oo)\\_______
            (__)\\       )\\/\\
                ||----w |
                ||     ||
        `}
      </Typography>
    </Box>
  )
}

export const matrixCommand = (): React.ReactNode => {
  const matrixChars = "01アイウエオカキクケコサシスセソタチツテト"
  const lines = Array.from({ length: 10 }, () =>
    Array.from({ length: 60 }, () => matrixChars[Math.floor(Math.random() * matrixChars.length)]).join(""),
  )

  return (
    <Box>
      {lines.map((line, index) => (
        <Typography
          key={index}
          sx={{
            fontFamily: "monospace",
            fontSize: "12px",
            color: "#00ff00",
            lineHeight: 1.2,
          }}
        >
          {line}
        </Typography>
      ))}
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00aaff", mt: 1 }}>
        Wake up, Neo... The Matrix has you...
      </Typography>
    </Box>
  )
}

export const rickrollCommand = (): React.ReactNode => {
  return (
    <Box>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ff00ff", mb: 1 }}>
        ♪ Never gonna give you up ♪
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ff00ff", mb: 1 }}>
        ♪ Never gonna let you down ♪
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ff00ff", mb: 1 }}>
        ♪ Never gonna run around and desert you ♪
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00aaff", mt: 2 }}>
        You just got rickrolled! 🎵
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#888888", mt: 1 }}>
        (Opening YouTube would be too obvious...)
      </Typography>
    </Box>
  )
}

export const coffeeCommand = (): React.ReactNode => {
  return (
    <Box>
      <Typography component="pre" sx={{ fontFamily: "monospace", fontSize: "14px", color: "#8B4513" }}>
        {`
        (  )   (   )  )
         ) (   )  (  (
         ( )  (    ) )
         _____________
        <_____________> ___
        |             |/ _ \\
        |               | | |
        |               |_| |
     ___|             |\\___/
    /    \\___________/    \\
    \\_____________________/
        `}
      </Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00aaff", mt: 1 }}>
        ☕ Coffee break! Refueling developer energy...
      </Typography>
    </Box>
  )
}
