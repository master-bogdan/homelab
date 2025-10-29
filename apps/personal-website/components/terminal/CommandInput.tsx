"use client"

import type React from "react"
import { useState, useRef, useEffect } from "react"
import { Box } from "@mui/material"
import { Prompt } from "./Prompt"

interface CommandInputProps {
  onCommand: (command: string) => void
  disabled?: boolean
}

export const CommandInput: React.FC<CommandInputProps> = ({ onCommand, disabled = false }) => {
  const [input, setInput] = useState("")
  const [cursorVisible, setCursorVisible] = useState(true)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (!disabled && inputRef.current) {
      inputRef.current.focus()
    }
  }, [disabled])

  useEffect(() => {
    const interval = setInterval(() => {
      setCursorVisible((prev) => !prev)
    }, 530)
    return () => clearInterval(interval)
  }, [])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (input.trim()) {
      onCommand(input.trim())
      setInput("")
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleSubmit(e)
    }
  }

  return (
    <Box
      component="form"
      onSubmit={handleSubmit}
      sx={{
        display: "flex",
        alignItems: "center",
        gap: 1,
      }}
    >
      <Prompt showCursor={false} />
      <Box
        sx={{
          position: "relative",
          flex: 1,
          fontFamily: "monospace",
          fontSize: "14px",
          color: "#00ff00",
          display: "flex",
          alignItems: "center",
        }}
      >
        <input
          ref={inputRef}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          disabled={disabled}
          autoComplete="off"
          spellCheck={false}
          style={{
            position: "absolute",
            left: 0,
            top: 0,
            width: "100%",
            height: "100%",
            opacity: 0,
            cursor: "default",
            border: "none",
            outline: "none",
            background: "transparent",
          }}
        />
        <Box
          component="span"
          sx={{
            fontFamily: "monospace",
            fontSize: "14px",
            color: "#00ff00",
            whiteSpace: "pre",
          }}
        >
          {input}
          <Box
            component="span"
            sx={{
              display: "inline-block",
              width: "10px",
              height: "18px",
              backgroundColor: cursorVisible ? "#00ff00" : "transparent",
              marginLeft: "2px",
              verticalAlign: "middle",
            }}
          />
        </Box>
      </Box>
    </Box>
  )
}
