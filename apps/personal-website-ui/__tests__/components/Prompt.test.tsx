import { describe, it, expect } from "vitest"
import { render, screen } from "@testing-library/react"
import { Prompt } from "@/components/terminal/Prompt"

describe("Prompt", () => {
  it("renders the prompt with user and host", () => {
    render(<Prompt />)
    expect(screen.getByText(/bogdan@homelab/i)).toBeInTheDocument()
  })

  it("shows cursor when showCursor is true", () => {
    const { container } = render(<Prompt showCursor={true} />)
    expect(container.textContent).toContain("▊")
  })

  it("hides cursor when showCursor is false", () => {
    const { container } = render(<Prompt showCursor={false} />)
    expect(container.textContent).not.toContain("▊")
  })
})
