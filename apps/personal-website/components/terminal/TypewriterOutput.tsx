"use client"

import type React from "react"
import { useState, useEffect, cloneElement, isValidElement } from "react"
import { Box } from "@mui/material"

interface TypewriterOutputProps {
  children: React.ReactNode
  speed?: number
  onComplete?: () => void
}

export const TypewriterOutput: React.FC<TypewriterOutputProps> = ({ children, speed = 15, onComplete }) => {
  const [charCount, setCharCount] = useState(0)
  const [totalChars, setTotalChars] = useState(0)
  const [isComplete, setIsComplete] = useState(false)

  useEffect(() => {
    // Count total characters in the tree
    const countChars = (node: React.ReactNode): number => {
      if (typeof node === "string") return node.length
      if (typeof node === "number") return String(node).length
      if (node === null || node === undefined) return 0

      if (Array.isArray(node)) {
        return node.reduce((sum, child) => sum + countChars(child), 0)
      }

      if (isValidElement(node)) {
        return countChars(node.props.children)
      }

      return 0
    }

    const total = countChars(children)
    setTotalChars(total)
    setCharCount(0)
    setIsComplete(false)

    let current = 0
    const timer = setInterval(() => {
      if (current < total) {
        current++
        setCharCount(current)
      } else {
        setIsComplete(true)
        clearInterval(timer)
        if (onComplete) {
          onComplete()
        }
      }
    }, speed)

    return () => clearInterval(timer)
  }, [children, speed, onComplete])

  // Recursively slice content while preserving structure
  const sliceContent = (node: React.ReactNode, remainingChars: number): [React.ReactNode, number] => {
    if (remainingChars <= 0) return [null, 0]

    if (typeof node === "string") {
      const sliced = node.slice(0, remainingChars)
      return [sliced, remainingChars - sliced.length]
    }

    if (typeof node === "number") {
      const str = String(node)
      const sliced = str.slice(0, remainingChars)
      return [sliced, remainingChars - sliced.length]
    }

    if (node === null || node === undefined) {
      return [null, remainingChars]
    }

    if (Array.isArray(node)) {
      const result: React.ReactNode[] = []
      let remaining = remainingChars

      for (const child of node) {
        if (remaining <= 0) break
        const [slicedChild, newRemaining] = sliceContent(child, remaining)
        if (slicedChild !== null) {
          result.push(slicedChild)
        }
        remaining = newRemaining
      }

      return [result, remaining]
    }

    if (isValidElement(node)) {
      const [slicedChildren, remaining] = sliceContent(node.props.children, remainingChars)

      if (slicedChildren === null) return [null, remainingChars]

      return [cloneElement(node, { ...node.props }, slicedChildren), remaining]
    }

    return [null, remainingChars]
  }

  if (isComplete) {
    return <Box>{children}</Box>
  }

  const [slicedContent] = sliceContent(children, charCount)

  return <Box>{slicedContent}</Box>
}
