"use client"

import { useState, useEffect } from "react"

interface UseTypingAnimationProps {
  text: string
  speed?: number
  onComplete?: () => void
}

export const useTypingAnimation = ({ text, speed = 50, onComplete }: UseTypingAnimationProps) => {
  const [displayedText, setDisplayedText] = useState("")
  const [currentIndex, setCurrentIndex] = useState(0)

  useEffect(() => {
    if (currentIndex < text.length) {
      const timeout = setTimeout(() => {
        setDisplayedText((prev) => prev + text[currentIndex])
        setCurrentIndex((prev) => prev + 1)
      }, speed)

      return () => clearTimeout(timeout)
    } else if (onComplete) {
      onComplete()
    }
  }, [currentIndex, text, speed, onComplete])

  return displayedText
}
