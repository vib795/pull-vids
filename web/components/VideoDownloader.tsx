'use client'

import { useState } from 'react'

export default function VideoDownloader() {
  const [url, setUrl] = useState('')
  const [quality, setQuality] = useState('best')
  const [audioOnly, setAudioOnly] = useState(false)
  const [showInstructions, setShowInstructions] = useState(false)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setShowInstructions(true)
  }

  const cliCommand = audioOnly
    ? `pull-vids --audio '${url}'`
    : quality !== 'best'
    ? `pull-vids --quality ${quality} '${url}'`
    : `pull-vids '${url}'`

  return (
    <div className="max-w-3xl mx-auto">
      {/* Notice Banner */}
      <div className="mb-8 bg-yellow-50 dark:bg-yellow-900/20 border-2 border-yellow-200 dark:border-yellow-800 rounded-xl p-6">
        <div className="flex items-start gap-3">
          <span className="text-2xl">⚠️</span>
          <div>
            <h3 className="font-bold text-lg mb-2">Web Downloads Currently Unavailable</h3>
            <p className="text-sm text-gray-700 dark:text-gray-300 mb-3">
              Due to platform restrictions on serverless functions, the web interface cannot directly download videos.
            </p>
            <p className="text-sm text-gray-700 dark:text-gray-300">
              <strong>Use the CLI version instead</strong> - it's more powerful, supports 1000+ platforms, and works reliably!
            </p>
          </div>
        </div>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl p-8">
        <h2 className="text-2xl font-bold mb-6 text-center">
          Generate Download Command
        </h2>

        <form onSubmit={handleSubmit} className="space-y-6">
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
              Paste a video URL from any supported platform
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

          {/* Submit Button */}
          <button
            type="submit"
            disabled={!url}
            className="w-full bg-gradient-to-r from-blue-600 to-purple-600 text-white font-bold py-4 px-6 rounded-lg hover:from-blue-700 hover:to-purple-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-105 active:scale-95"
          >
            📋 Generate CLI Command
          </button>
        </form>

        {/* Instructions */}
        {showInstructions && url && (
          <div className="mt-8 pt-8 border-t border-gray-200 dark:border-gray-700 space-y-6">
            {/* Step 1: Install */}
            <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-6">
              <h3 className="font-bold text-lg mb-3 flex items-center gap-2">
                <span className="bg-blue-600 text-white rounded-full w-8 h-8 flex items-center justify-center text-sm">1</span>
                Install pull-vids CLI
              </h3>

              <div className="space-y-3">
                <div>
                  <p className="text-sm font-medium mb-2">macOS (Homebrew):</p>
                  <div className="bg-gray-900 text-green-400 p-3 rounded font-mono text-sm overflow-x-auto">
                    brew tap vib795/tap && brew install pull-vids
                  </div>
                </div>

                <div>
                  <p className="text-sm font-medium mb-2">Other platforms:</p>
                  <a
                    href="https://github.com/vib795/pull-vids#installation"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-blue-600 dark:text-blue-400 hover:underline text-sm"
                  >
                    View all installation options →
                  </a>
                </div>
              </div>
            </div>

            {/* Step 2: Run Command */}
            <div className="bg-green-50 dark:bg-green-900/20 rounded-lg p-6">
              <h3 className="font-bold text-lg mb-3 flex items-center gap-2">
                <span className="bg-green-600 text-white rounded-full w-8 h-8 flex items-center justify-center text-sm">2</span>
                Run this command
              </h3>

              <div className="bg-gray-900 text-green-400 p-4 rounded font-mono text-sm overflow-x-auto relative">
                <code>{cliCommand}</code>
                <button
                  onClick={() => navigator.clipboard.writeText(cliCommand)}
                  className="absolute top-2 right-2 bg-gray-700 hover:bg-gray-600 text-white px-3 py-1 rounded text-xs transition"
                >
                  Copy
                </button>
              </div>

              <p className="text-xs text-gray-600 dark:text-gray-400 mt-2">
                Click to copy, then paste in your terminal
              </p>
            </div>

            {/* Why CLI? */}
            <div className="bg-purple-50 dark:bg-purple-900/20 rounded-lg p-6">
              <h3 className="font-bold text-lg mb-3">✨ Why use CLI?</h3>
              <ul className="space-y-2 text-sm">
                <li className="flex items-start gap-2">
                  <span className="text-green-600">✅</span>
                  <span>Supports <strong>1000+ platforms</strong> (not just YouTube)</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-green-600">✅</span>
                  <span>No file size limits</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-green-600">✅</span>
                  <span>Playlist downloads</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-green-600">✅</span>
                  <span>Faster and more reliable</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-green-600">✅</span>
                  <span>Works offline</span>
                </li>
              </ul>
            </div>
          </div>
        )}

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
