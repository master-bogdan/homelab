import type { Command } from "./types"
import { whoamiCommand } from "./commands/whoami"
import { projectsCommand } from "./commands/projects"
import { contactsCommand } from "./commands/contacts"
import { resumeCommand } from "./commands/resume"
import { helpCommand } from "./commands/help"
import { blogCommand } from "./commands/blog"
import {
  neofetchCommand,
  sudoCommand,
  hackCommand,
  fortuneCommand,
  cowsayCommand,
  matrixCommand,
  rickrollCommand,
  coffeeCommand,
} from "./commands/hidden"

export const commandRegistry: Record<string, Command> = {
  whoami: {
    name: "whoami",
    description: "Display information about me",
    execute: whoamiCommand,
  },
  projects: {
    name: "projects",
    description: "List featured projects",
    execute: projectsCommand,
  },
  blog: {
    name: "blog",
    description: "Read blog posts",
    execute: blogCommand,
  },
  contacts: {
    name: "contacts",
    description: "Get contact information",
    execute: contactsCommand,
  },
  resume: {
    name: "resume",
    description: "Download my resume",
    execute: resumeCommand,
  },
  clear: {
    name: "clear",
    description: "Clear the terminal output",
    execute: () => null, // Handled in Terminal component
  },
  help: {
    name: "help",
    description: "Show available commands",
    execute: helpCommand,
  },
  neofetch: {
    name: "neofetch",
    description: "Display system information",
    execute: neofetchCommand,
    hidden: true,
  },
  sudo: {
    name: "sudo",
    description: "Execute command as superuser",
    execute: sudoCommand,
    hidden: true,
  },
  hack: {
    name: "hack",
    description: "Initiate hack sequence",
    execute: hackCommand,
    hidden: true,
  },
  fortune: {
    name: "fortune",
    description: "Get your fortune",
    execute: fortuneCommand,
    hidden: true,
  },
  cowsay: {
    name: "cowsay",
    description: "Make a cow say something",
    execute: cowsayCommand,
    hidden: true,
  },
  matrix: {
    name: "matrix",
    description: "Enter the Matrix",
    execute: matrixCommand,
    hidden: true,
  },
  rickroll: {
    name: "rickroll",
    description: "Never gonna give you up",
    execute: rickrollCommand,
    hidden: true,
  },
  coffee: {
    name: "coffee",
    description: "Get some coffee",
    execute: coffeeCommand,
    hidden: true,
  },
}

export const getVisibleCommands = (): string[] => {
  return Object.values(commandRegistry)
    .filter((cmd) => !cmd.hidden)
    .map((cmd) => cmd.name)
}
