import type React from "react"
import { Box, Typography, Link } from "@mui/material"

export const projectsCommand = (): React.ReactNode => {
  const projects = [
    {
      name: "Cloud Infrastructure Platform",
      tech: "Kubernetes, AWS, Terraform",
      description: "Scalable multi-tenant platform serving 10k+ users",
      url: "https://github.com/bogdan", // Added URL field
    },
    {
      name: "Real-time Analytics Engine",
      tech: "Node.js, TypeScript, Redis",
      description: "High-throughput data processing pipeline",
      url: "https://github.com/bogdan",
    },
    {
      name: "Microservices Architecture",
      tech: "Go, gRPC, Docker",
      description: "Distributed system with 99.9% uptime",
      url: "https://github.com/bogdan",
    },
    {
      name: "CI/CD Pipeline Automation",
      tech: "GitHub Actions, ArgoCD",
      description: "Zero-downtime deployment system",
      url: "https://github.com/bogdan",
    },
  ]

  return (
    <Box>
      <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00aaff", mb: 1 }}>
        Featured Projects:
      </Typography>
      {projects.map((project, index) => (
        <Box key={index} sx={{ mb: 2 }}>
          <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#00ff00" }}>
            {index + 1}.{" "}
            <Link
              href={project.url}
              target="_blank"
              rel="noopener noreferrer"
              sx={{
                color: "#00ff00",
                textDecoration: "none",
                "&:hover": {
                  textDecoration: "underline",
                  color: "#00aaff",
                },
              }}
            >
              {project.name}
            </Link>
          </Typography>
          <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#888888", ml: 2 }}>
            Tech: {project.tech}
          </Typography>
          <Typography sx={{ fontFamily: "monospace", fontSize: "14px", color: "#ffffff", ml: 2 }}>
            {project.description}
          </Typography>
        </Box>
      ))}
    </Box>
  )
}
