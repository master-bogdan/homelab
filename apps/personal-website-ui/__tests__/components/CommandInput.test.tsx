import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { CommandInput } from "@/components/terminal/CommandInput"

describe("CommandInput", () => {
  it("renders input field", () => {
    const onCommand = vi.fn()
    render(<CommandInput onCommand={onCommand} />)
    const input = screen.getByRole("textbox")
    expect(input).toBeInTheDocument()
  })

  it("calls onCommand when form is submitted", () => {
    const onCommand = vi.fn()
    render(<CommandInput onCommand={onCommand} />)
    const input = screen.getByRole("textbox")

    fireEvent.change(input, { target: { value: "whoami" } })
    fireEvent.submit(input.closest("form")!)

    expect(onCommand).toHaveBeenCalledWith("whoami")
  })

  it("clears input after submission", () => {
    const onCommand = vi.fn()
    render(<CommandInput onCommand={onCommand} />)
    const input = screen.getByRole("textbox") as HTMLInputElement

    fireEvent.change(input, { target: { value: "help" } })
    fireEvent.submit(input.closest("form")!)

    expect(input.value).toBe("")
  })

  it("does not submit empty commands", () => {
    const onCommand = vi.fn()
    render(<CommandInput onCommand={onCommand} />)
    const input = screen.getByRole("textbox")

    fireEvent.change(input, { target: { value: "   " } })
    fireEvent.submit(input.closest("form")!)

    expect(onCommand).not.toHaveBeenCalled()
  })
})
