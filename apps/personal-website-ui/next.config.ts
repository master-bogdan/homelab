// next.config.ts
import type { NextConfig } from "next"

const nextConfig: NextConfig = {
  reactStrictMode: true,
  poweredByHeader: false,
  compress: true,
  reactCompiler: true,
  output: "export",
  images: {
    unoptimized: true,
  },
}

export default nextConfig
