import { describe, it, expect } from "vitest"
import { render, screen } from "@testing-library/react"
import { Banner } from "@/components/terminal/Banner"

describe("Banner", () => {
  it("renders the banner image", () => {
    render(<Banner />)
    const image = screen.getByRole("img", { name: /bogdan shchavinskyi banner/i })
    expect(image).toHaveAttribute("src", expect.stringContaining("/banner.svg"))
  })
})
