import type React from "react"
export interface CommandOutput {
  command: string
  output: React.ReactNode
  timestamp: Date
}

export interface Command {
  name: string
  description: string
  execute: () => React.ReactNode
  hidden?: boolean
}
