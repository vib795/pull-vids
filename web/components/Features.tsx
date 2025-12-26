export default function Features() {
  const features = [
    {
      icon: '🌐',
      title: '1000+ Websites',
      description: 'Download from YouTube, Vimeo, Twitter, TikTok, Instagram, and 1000+ more platforms',
    },
    {
      icon: '🆓',
      title: '100% Free',
      description: 'No subscriptions, no trials, no hidden costs. Completely free forever.',
    },
    {
      icon: '🎥',
      title: 'Multiple Qualities',
      description: 'Choose from 360p to 4K. Download in the quality that fits your needs.',
    },
    {
      icon: '🎵',
      title: 'Audio Extraction',
      description: 'Extract audio as MP3, M4A, or other formats. Perfect for music and podcasts.',
    },
    {
      icon: '📱',
      title: 'Easy to Use',
      description: 'Simple web interface. Just paste the URL and click download. No technical knowledge required.',
    },
    {
      icon: '🔒',
      title: 'Privacy Focused',
      description: 'No tracking, no data collection. Your downloads are completely private.',
    },
  ]

  return (
    <div id="features" className="mt-24">
      <h2 className="text-4xl font-bold text-center mb-12">Why pull-vids?</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {features.map((feature, index) => (
          <div
            key={index}
            className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-lg hover:shadow-xl transition-shadow"
          >
            <div className="text-4xl mb-4">{feature.icon}</div>
            <h3 className="text-xl font-bold mb-2">{feature.title}</h3>
            <p className="text-gray-600 dark:text-gray-400">{feature.description}</p>
          </div>
        ))}
      </div>

      {/* Limitations */}
      <div className="mt-16 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-xl p-6">
        <h3 className="text-xl font-bold mb-3 flex items-center gap-2">
          ⚠️ Important Limitations
        </h3>
        <p className="text-gray-700 dark:text-gray-300 mb-3">
          This tool <strong>cannot</strong> download DRM-protected content from:
        </p>
        <ul className="list-disc list-inside text-gray-700 dark:text-gray-300 space-y-1">
          <li>Netflix, Amazon Prime Video, Disney+, Hulu, HBO Max</li>
          <li>Apple TV+, Paramount+, Peacock, or other paid streaming services</li>
        </ul>
        <p className="text-gray-600 dark:text-gray-400 text-sm mt-3">
          Only download content you have the legal right to download. Respect copyright laws and platform Terms of Service.
        </p>
      </div>
    </div>
  )
}
