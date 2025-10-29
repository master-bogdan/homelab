"use client"

import type React from "react"
import { Box, Typography, Link } from "@mui/material"
import ReactMarkdown from "react-markdown"
import { getBlogPosts, getBlogPost } from "../blog"

export const blogCommand = (args?: string[]): React.ReactNode => {
  // If no args, list all posts
  if (!args || args.length === 0) {
    const posts = getBlogPosts()
    return (
      <Box>
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00aaff", mb: 1 }}>Blog Posts:</Typography>
        {posts.map((post, index) => (
          <Box key={post.slug} sx={{ mb: 1 }}>
            <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00ff00" }}>
              {index + 1}.{" "}
              <Link
                key={post.slug}
                component="button"
                onClick={() => {
                  // Trigger blog command with slug
                  const event = new CustomEvent("terminal-command", { detail: `blog ${post.slug}` })
                  window.dispatchEvent(event)
                }}
                sx={{
                  color: "#00ff00",
                  textDecoration: "none",
                  cursor: "pointer",
                  background: "none",
                  border: "none",
                  fontFamily: "monospace",
                  fontSize: "14px",
                  padding: 0,
                  "&:hover": {
                    textDecoration: "underline",
                    color: "#00aaff",
                  },
                }}
              >
                {post.title}
              </Link>
            </Typography>
            <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#888888", ml: 2 }}>
              Date: {post.date}
            </Typography>
          </Box>
        ))}
        <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#888888", mt: 2 }}>
          Click on a post title to read it, or use: blog [slug]
        </Typography>
      </Box>
    )
  }

  // If args provided, show specific post
  const slug = args[0]
  const post = getBlogPost(slug)

  if (!post) {
    return (
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ff0000" }}>
        Error: Post &apos;{slug}&apos; not found. Use &apos;blog&apos; to list all posts.
      </Typography>
    )
  }

  return (
    <Box>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00aaff", mb: 1 }}>{post.title}</Typography>
      <Typography sx={{ fontFamily: "monospace", fontSize: "12px", color: "#888888", mb: 2 }}>
        Published: {post.date}
      </Typography>
      <Box
        sx={{
          "& h1": {
            fontFamily: "monospace",
            fontSize: "18px",
            color: "#00aaff",
            mb: 1,
            mt: 2,
          },
          "& h2": {
            fontFamily: "monospace",
            fontSize: "16px",
            color: "#00ff00",
            mb: 1,
            mt: 1.5,
          },
          "& p": {
            fontFamily: "monospace",
            fontSize: "14px",
            color: "#ffffff",
            mb: 1,
            lineHeight: 1.6,
          },
          "& ul, & ol": {
            fontFamily: "monospace",
            fontSize: "14px",
            color: "#ffffff",
            ml: 2,
            mb: 1,
          },
          "& li": {
            mb: 0.5,
          },
          "& code": {
            fontFamily: "monospace",
            fontSize: "13px",
            color: "#00ff00",
            backgroundColor: "#1a1a1a",
            padding: "2px 6px",
            borderRadius: "3px",
          },
          "& pre": {
            fontFamily: "monospace",
            fontSize: "13px",
            color: "#00ff00",
            backgroundColor: "#1a1a1a",
            padding: "12px",
            borderRadius: "4px",
            overflow: "auto",
            mb: 1,
          },
          "& pre code": {
            backgroundColor: "transparent",
            padding: 0,
          },
          "& a": {
            color: "#00aaff",
            textDecoration: "none",
            "&:hover": {
              textDecoration: "underline",
            },
          },
        }}
      >
        <ReactMarkdown>{post.content}</ReactMarkdown>
      </Box>
    </Box>
  )
}
