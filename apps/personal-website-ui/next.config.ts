// next.config.ts
import type { NextConfig } from "next"

const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? ""
const assetPrefix = basePath ? `${basePath}/` : undefined

const nextConfig: NextConfig = {
  reactStrictMode: true,
  poweredByHeader: false,
  compress: true,
  reactCompiler: true,
  output: "export",
  basePath,
  assetPrefix,
  trailingSlash: true,
  images: {
    unoptimized: true,
  },
  distDir: "build"
}

export default nextConfig
