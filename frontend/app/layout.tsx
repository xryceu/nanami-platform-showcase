import type { Metadata, Viewport } from "next";
import type { ReactNode } from "react";

import { ShowcaseNav } from "@/components/showcase-nav";

import "./globals.css";

export const metadata: Metadata = {
  title: "Nanami runtime excerpt",
  description: "A public-safe excerpt of the Nanami Client App.",
};

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  viewportFit: "cover",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>
        <ShowcaseNav />
        {children}
      </body>
    </html>
  );
}
