import { describe, it, expect } from "vitest"
import { getBlogPosts, getBlogPost } from "@/lib/blog"

describe("Blog", () => {
  it("returns all blog posts", () => {
    const posts = getBlogPosts()
    expect(posts.length).toBeGreaterThan(0)
    expect(posts[0]).toHaveProperty("slug")
    expect(posts[0]).toHaveProperty("title")
    expect(posts[0]).toHaveProperty("date")
    expect(posts[0]).toHaveProperty("content")
  })

  it("returns a specific blog post by slug", () => {
    const post = getBlogPost("welcome")
    expect(post).toBeDefined()
    expect(post?.slug).toBe("welcome")
    expect(post?.title).toContain("Welcome")
  })

  it("returns undefined for non-existent slug", () => {
    const post = getBlogPost("non-existent-post")
    expect(post).toBeUndefined()
  })
})
