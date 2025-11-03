import { describe, it, expect } from "vitest"
import { commandRegistry, getVisibleCommands } from "@/lib/commandRegistry"

describe("Command Registry", () => {
  it("contains all required commands", () => {
    const requiredCommands = ["whoami", "projects", "blog", "contacts", "resume", "help"]
    requiredCommands.forEach((cmd) => {
      expect(commandRegistry[cmd]).toBeDefined()
    })
  })

  it("contains hidden commands", () => {
    const hiddenCommands = ["motd", "fortune", "skills", "theme"]
    hiddenCommands.forEach((cmd) => {
      expect(commandRegistry[cmd]).toBeDefined()
      expect(commandRegistry[cmd].hidden).toBe(true)
    })
  })

  it("getVisibleCommands returns only non-hidden commands", () => {
    const visibleCommands = getVisibleCommands()
    expect(visibleCommands).toContain("whoami")
    expect(visibleCommands).toContain("help")
    expect(visibleCommands).not.toContain("motd")
    expect(visibleCommands).not.toContain("fortune")
  })

  it("all commands have required properties", () => {
    Object.values(commandRegistry).forEach((cmd) => {
      expect(cmd).toHaveProperty("name")
      expect(cmd).toHaveProperty("description")
      expect(cmd).toHaveProperty("execute")
      expect(typeof cmd.execute).toBe("function")
    })
  })
})
