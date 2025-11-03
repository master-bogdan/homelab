import { describe, it, expect } from "vitest"
import { render, screen } from "@testing-library/react"
import { Banner } from "@/components/terminal/Banner"

describe("Banner", () => {
  it("renders the ASCII banner", () => {
    render(<Banner />)
    expect(screen.getByText(/BOGDAN/i)).toBeInTheDocument()
  })

  it("displays professional title", () => {
    const { container } = render(<Banner />)
    expect(container.textContent).toContain("Senior Software Engineer")
  })

  it("displays technical specializations", () => {
    const { container } = render(<Banner />)
    expect(container.textContent).toContain("Backend")
    expect(container.textContent).toContain("Platform")
    expect(container.textContent).toContain("DevOps")
  })
})
