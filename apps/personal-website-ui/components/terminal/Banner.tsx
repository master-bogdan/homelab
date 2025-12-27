import type React from "react"
import Image from 'next/image';

const basePath = (process.env.NEXT_PUBLIC_BASE_PATH ?? "").replace(/\/$/, "")
const src = `${basePath}/banner.svg`

export const Banner: React.FC = () => {
  return (
    <Image
      priority
      src={src}
      height={32}
      width={32}
      alt="Follow us on Twitter"
    />
  )
}
