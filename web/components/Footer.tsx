export default function Footer() {
  return (
    <footer className="bg-gray-100 dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700 mt-16">
      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* About */}
          <div>
            <h3 className="font-bold text-lg mb-4">About pull-vids</h3>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              Free, open-source video downloader supporting 1000+ websites.
              Built with ❤️ for the community.
            </p>
          </div>

          {/* Links */}
          <div>
            <h3 className="font-bold text-lg mb-4">Links</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <a
                  href="https://github.com/vib795/pull-vids"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition"
                >
                  GitHub Repository
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/vib795/pull-vids/issues"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition"
                >
                  Report an Issue
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/vib795/homebrew-tap"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition"
                >
                  Homebrew Tap
                </a>
              </li>
            </ul>
          </div>

          {/* Disclaimer */}
          <div>
            <h3 className="font-bold text-lg mb-4">Disclaimer</h3>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              This tool is for personal use only. Please respect copyright laws
              and platform Terms of Service. Only download content you have the
              right to download.
            </p>
          </div>
        </div>

        <div className="border-t border-gray-200 dark:border-gray-700 mt-8 pt-8 text-center">
          <p className="text-gray-600 dark:text-gray-400 text-sm">
            © {new Date().getFullYear()} pull-vids. Open source under MIT License.
          </p>
          <p className="text-gray-500 dark:text-gray-500 text-xs mt-2">
            Not affiliated with YouTube, Vimeo, Twitter, TikTok, Instagram, or any video platform.
          </p>
        </div>
      </div>
    </footer>
  )
}
