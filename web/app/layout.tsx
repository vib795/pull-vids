import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'pull-vids - Free Video Downloader',
  description: 'Download videos from 1000+ websites for free. YouTube, Vimeo, Twitter, TikTok, Instagram, and more!',
  keywords: ['video downloader', 'youtube downloader', 'free', 'vimeo', 'tiktok', 'instagram'],
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>{children}</body>
    </html>
  )
}
