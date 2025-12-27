"use client"

import type React from "react"
import { useState, useEffect } from "react"
import { Box } from "@mui/material"
import { BOOT_MESSAGES } from "@/lib/constants"
import { TypewriterOutput } from "./TypewriterOutput"

interface BootSequenceProps {
  onComplete: () => void
}

export const BootSequence: React.FC<BootSequenceProps> = ({ onComplete }) => {
  const [visibleMessages, setVisibleMessages] = useState<string[]>([])
  const [currentIndex, setCurrentIndex] = useState(0)
  const [isTyping, setIsTyping] = useState(false)

  useEffect(() => {
    if (currentIndex < BOOT_MESSAGES.length && !isTyping) {
      setIsTyping(true)
      setVisibleMessages((prev) => [...prev, BOOT_MESSAGES[currentIndex]])
    } else if (currentIndex >= BOOT_MESSAGES.length) {
      const timeout = setTimeout(onComplete, 300)
      return () => clearTimeout(timeout)
    }
  }, [currentIndex, isTyping, onComplete])

  const handleMessageComplete = (index: number) => {
    if (index === visibleMessages.length - 1) {
      setIsTyping(false)
      setCurrentIndex((prev) => prev + 1)
    }
  }

  return (
    <Box>
      {visibleMessages.map((message, index) => (
        <Box
          key={index}
          sx={{
            fontFamily: "var(--font-mono)",
            fontSize: "14px",
            color: "#00ff00",
            mb: 0.5,
          }}
        >
          {index === visibleMessages.length - 1 ? (
            <TypewriterOutput speed={20} onComplete={() => handleMessageComplete(index)}>
              [ OK ] {message}
            </TypewriterOutput>
          ) : (
            <>[ OK ] {message}</>
          )}
        </Box>
      ))}
    </Box>
  )
}
