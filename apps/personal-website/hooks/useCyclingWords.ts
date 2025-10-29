"use client"

import { useState, useEffect } from "react"

export const useCyclingWords = (words: string[], interval = 2000) => {
  const [currentIndex, setCurrentIndex] = useState(0)

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentIndex((prev) => (prev + 1) % words.length)
    }, interval)

    return () => clearInterval(timer)
  }, [words, interval])

  return words[currentIndex]
}
