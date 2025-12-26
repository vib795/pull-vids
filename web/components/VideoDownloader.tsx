'use client'

import { useState } from 'react'

export default function VideoDownloader() {
  const [url, setUrl] = useState('')
  const [quality, setQuality] = useState('best')
  const [audioOnly, setAudioOnly] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

  const handleDownload = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    setSuccess(false)

    try {
      const response = await fetch('/api/download', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          url,
          quality,
          audioOnly,
        }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.error || 'Download failed')
      }

      // Get the download link from the response
      const data = await response.json()

      if (data.downloadUrl) {
        // Trigger download
        window.location.href = data.downloadUrl
        setSuccess(true)
        setUrl('')
      } else {
        throw new Error('No download URL received')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-3xl mx-auto">
      <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl p-8">
        <form onSubmit={handleDownload} className="space-y-6">
          {/* URL Input */}
          <div>
            <label htmlFor="url" className="block text-sm font-medium mb-2">
              Video URL
            </label>
            <input
              id="url"
              type="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="https://www.youtube.com/watch?v=..."
              required
              className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
            />
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Paste a video URL from YouTube, Vimeo, Twitter, TikTok, Instagram, or any supported platform
            </p>
          </div>

          {/* Options */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* Quality Selection */}
            <div>
              <label htmlFor="quality" className="block text-sm font-medium mb-2">
                Quality
              </label>
              <select
                id="quality"
                value={quality}
                onChange={(e) => setQuality(e.target.value)}
                disabled={audioOnly}
                className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-700 dark:text-white disabled:opacity-50"
              >
                <option value="best">Best Quality</option>
                <option value="2160p">4K (2160p)</option>
                <option value="1440p">2K (1440p)</option>
                <option value="1080p">Full HD (1080p)</option>
                <option value="720p">HD (720p)</option>
                <option value="480p">480p</option>
                <option value="360p">360p</option>
              </select>
            </div>

            {/* Audio Only Toggle */}
            <div className="flex items-end">
              <label className="flex items-center cursor-pointer w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition">
                <input
                  type="checkbox"
                  checked={audioOnly}
                  onChange={(e) => setAudioOnly(e.target.checked)}
                  className="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
                />
                <span className="ml-3 text-sm font-medium">
                  🎵 Audio Only (MP3)
                </span>
              </label>
            </div>
          </div>

          {/* Error Message */}
          {error && (
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
              <p className="text-red-800 dark:text-red-300 text-sm">
                ❌ {error}
              </p>
            </div>
          )}

          {/* Success Message */}
          {success && (
            <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4">
              <p className="text-green-800 dark:text-green-300 text-sm">
                ✅ Download started! Check your downloads folder.
              </p>
            </div>
          )}

          {/* Submit Button */}
          <button
            type="submit"
            disabled={loading || !url}
            className="w-full bg-gradient-to-r from-blue-600 to-purple-600 text-white font-bold py-4 px-6 rounded-lg hover:from-blue-700 hover:to-purple-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-105 active:scale-95"
          >
            {loading ? (
              <span className="flex items-center justify-center gap-2">
                <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                    fill="none"
                  />
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  />
                </svg>
                Downloading...
              </span>
            ) : (
              '⬇️ Download Video'
            )}
          </button>
        </form>

        {/* Quick Examples */}
        <div className="mt-8 pt-8 border-t border-gray-200 dark:border-gray-700">
          <h3 className="text-sm font-medium mb-3 text-gray-700 dark:text-gray-300">
            Try with these examples:
          </h3>
          <div className="flex flex-wrap gap-2">
            {[
              { name: 'YouTube', url: 'https://www.youtube.com/watch?v=dQw4w9WgXcQ' },
              { name: 'Vimeo', url: 'https://vimeo.com/148751763' },
            ].map((example) => (
              <button
                key={example.name}
                type="button"
                onClick={() => setUrl(example.url)}
                className="px-3 py-1 text-xs bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-full hover:bg-gray-200 dark:hover:bg-gray-600 transition"
              >
                {example.name}
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
