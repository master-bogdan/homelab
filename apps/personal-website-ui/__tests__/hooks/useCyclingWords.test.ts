import { describe, it, expect, vi, beforeEach, afterEach } from "vitest"
import { renderHook, waitFor } from "@testing-library/react"
import { useCyclingWords } from "@/hooks/useCyclingWords"

describe("useCyclingWords", () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it("returns the first word initially", () => {
    const words = ["one", "two", "three"]
    const { result } = renderHook(() => useCyclingWords(words, 1000))
    expect(result.current).toBe("one")
  })

  it("cycles through words at specified interval", async () => {
    const words = ["one", "two", "three"]
    const { result } = renderHook(() => useCyclingWords(words, 1000))

    expect(result.current).toBe("one")

    vi.advanceTimersByTime(1000)
    await waitFor(() => expect(result.current).toBe("two"))

    vi.advanceTimersByTime(1000)
    await waitFor(() => expect(result.current).toBe("three"))

    vi.advanceTimersByTime(1000)
    await waitFor(() => expect(result.current).toBe("one"))
  })
})
