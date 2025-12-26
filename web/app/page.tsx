'use client'

import VideoDownloader from '@/components/VideoDownloader'
import Header from '@/components/Header'
import Footer from '@/components/Footer'
import Features from '@/components/Features'

export default function Home() {
  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900">
      <Header />

      <div className="container mx-auto px-4 py-16">
        {/* Hero Section */}
        <div className="text-center mb-16">
          <h1 className="text-6xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600">
            pull-vids
          </h1>
          <p className="text-2xl text-gray-700 dark:text-gray-300 mb-4">
            CLI Command Generator
          </p>
          <p className="text-xl text-gray-600 dark:text-gray-400 mb-8">
            Get the command to download videos from 1000+ websites
          </p>
          <div className="flex justify-center gap-4 text-sm text-gray-600 dark:text-gray-400">
            <span>✓ YouTube</span>
            <span>✓ Vimeo</span>
            <span>✓ Twitter</span>
            <span>✓ TikTok</span>
            <span>✓ Instagram</span>
            <span>✓ 1000+ more</span>
          </div>
        </div>

        {/* Download Section */}
        <VideoDownloader />

        {/* Features Section */}
        <Features />
      </div>

      <Footer />
    </main>
  )
}
