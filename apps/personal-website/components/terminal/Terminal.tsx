"use client"

import type React from "react"
import { useState, useEffect, useRef } from "react"
import { Box, Container } from "@mui/material"
import { BootSequence } from "./BootSequence"
import { Banner } from "./Banner"
import { CommandInput } from "./CommandInput"
import { QuickCommands } from "./QuickCommands"
import { TypewriterOutput } from "./TypewriterOutput"
import { commandRegistry, getVisibleCommands } from "@/lib/commandRegistry"
import type { CommandOutput } from "@/lib/types"

export const Terminal: React.FC = () => {
  const [isBooting, setIsBooting] = useState(true)
  const [showBanner, setShowBanner] = useState(false)
  const [currentOutput, setCurrentOutput] = useState<CommandOutput | null>(null)
  const outputRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (outputRef.current) {
      outputRef.current.scrollTop = outputRef.current.scrollHeight
    }
  }, [currentOutput])

  useEffect(() => {
    const handleCustomCommand = (event: Event) => {
      const customEvent = event as CustomEvent<string>
      if (customEvent.detail) {
        handleCommand(customEvent.detail)
      }
    }

    window.addEventListener("terminal-command", handleCustomCommand)
    return () => {
      window.removeEventListener("terminal-command", handleCustomCommand)
    }
  }, [])

  const handleBootComplete = () => {
    setIsBooting(false)
    setShowBanner(true)
    setTimeout(() => {
      const helpOutput = commandRegistry.help.execute([])
      setCurrentOutput({
        command: "help",
        output: helpOutput,
        timestamp: new Date(),
      })
    }, 500)
  }

  const handleCommand = (input: string) => {
    const [commandName, ...args] = input.toLowerCase().split(" ")

    if (commandName === "clear") {
      setCurrentOutput(null)
      // Show help after clearing
      setTimeout(() => {
        const helpOutput = commandRegistry.help.execute([])
        setCurrentOutput({
          command: "help",
          output: helpOutput,
          timestamp: new Date(),
        })
      }, 100)
      return
    }

    const command = commandRegistry[commandName]

    if (command) {
      const output = command.execute(args)
      setCurrentOutput({
        command: input,
        output,
        timestamp: new Date(),
      })
    } else {
      setCurrentOutput({
        command: input,
        output: (
          <Box sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ff0000" }}>
            Command not found: {commandName}
            <br />
            Type 'help' for available commands.
          </Box>
        ),
        timestamp: new Date(),
      })
    }
  }

  return (
    <Box
      sx={{
        minHeight: "100vh",
        backgroundColor: "#0a0a0a",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        padding: 2,
      }}
    >
      <Container
        maxWidth={false}
        sx={{
          maxWidth: "100ch",
          backgroundColor: "#0a0a0a",
          padding: 3,
          border: "1px solid #333",
          borderRadius: "4px",
          boxShadow: "0 0 20px rgba(0, 255, 0, 0.1)",
        }}
      >
        <Box
          ref={outputRef}
          sx={{
            minHeight: "60vh",
            maxHeight: "80vh",
            overflowY: "auto",
            mb: 2,
            "&::-webkit-scrollbar": {
              width: "8px",
            },
            "&::-webkit-scrollbar-track": {
              background: "#1a1a1a",
            },
            "&::-webkit-scrollbar-thumb": {
              background: "#00ff00",
              borderRadius: "4px",
            },
          }}
        >
          {isBooting && <BootSequence onComplete={handleBootComplete} />}
          {showBanner && (
            <>
              <Banner />
              <QuickCommands commands={getVisibleCommands()} onCommand={handleCommand} />
            </>
          )}
          {currentOutput && (
            <Box sx={{ mt: 2 }}>
              <Box sx={{ fontFamily: "monospace", fontSize: "14px", color: "#888888", mb: 1 }}>
                $ {currentOutput.command}
              </Box>
              <TypewriterOutput>{currentOutput.output}</TypewriterOutput>
            </Box>
          )}
        </Box>

        {!isBooting && <CommandInput onCommand={handleCommand} disabled={isBooting} />}
      </Container>
    </Box>
  )
}
