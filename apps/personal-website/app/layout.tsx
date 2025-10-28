export const metadata = {
  title: "Bogdan — retro terminal",
  description: "Senior Software Engineer — Go, TypeScript, DevOps",
};

import "../styles/globals.css";
import React from "react";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
