import Link from 'next/link'

export default function Header() {
  return (
    <header className="bg-white/80 dark:bg-gray-900/80 backdrop-blur-sm border-b border-gray-200 dark:border-gray-700 sticky top-0 z-50">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-2xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600">
              pull-vids
            </span>
            <span className="px-2 py-1 text-xs bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300 rounded-full">
              FREE
            </span>
          </div>

          <nav className="hidden md:flex items-center gap-6">
            <a
              href="#features"
              className="text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white transition"
            >
              Features
            </a>
            <a
              href="https://github.com/vib795/pull-vids"
              target="_blank"
              rel="noopener noreferrer"
              className="text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white transition"
            >
              GitHub
            </a>
            <a
              href="https://github.com/vib795/pull-vids#installation"
              target="_blank"
              rel="noopener noreferrer"
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
            >
              Install CLI
            </a>
          </nav>
        </div>
      </div>
    </header>
  )
}
