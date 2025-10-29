// next.config.ts
import type { NextConfig } from "next"

const nextConfig: NextConfig = {
  reactStrictMode: true,
  poweredByHeader: false,
  compress: true,
  reactCompiler: true,
  output: "standalone",
  images: {
    unoptimized: true,
  },
}

export default nextConfig
