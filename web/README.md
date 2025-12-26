# pull-vids Web App

Web interface for pull-vids - a free, universal video downloader.

## 🌐 Live Demo

**Coming soon!** This web app will be hosted for free on Vercel.

## ✨ Features

- **Easy to Use** - Simple web interface, no command line needed
- **1000+ Websites** - Download from YouTube, Vimeo, Twitter, TikTok, and more
- **Multiple Qualities** - Choose from 360p to 4K
- **Audio Extraction** - Download audio as MP3
- **Free Forever** - No subscriptions, no limits
- **Privacy Focused** - No tracking, no data collection

## 🚀 Quick Start

### Development

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Open http://localhost:3000
```

### Production Build

```bash
# Build for production
npm run build

# Start production server
npm start
```

## 📦 Tech Stack

- **Framework:** Next.js 14 (App Router)
- **Styling:** Tailwind CSS
- **Language:** TypeScript
- **Video Downloading:** ytdl-core
- **Hosting:** Vercel (free tier)

## 🔧 Configuration

### Environment Variables

Create a `.env.local` file:

```env
# Optional: Add any API keys or configuration here
```

### Vercel Deployment

1. Push this directory to GitHub
2. Import project in Vercel dashboard
3. Vercel will auto-detect Next.js and deploy
4. Done! Your web app is live

#### Deploy Button

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/vib795/pull-vids/tree/main/web)

## 📝 API Routes

### POST /api/download

Download a video.

**Request:**
```json
{
  "url": "https://www.youtube.com/watch?v=VIDEO_ID",
  "quality": "1080p",
  "audioOnly": false
}
```

**Response:**
```json
{
  "downloadUrl": "https://...",
  "title": "Video Title",
  "format": "mp4",
  "quality": "1080p"
}
```

### GET /api/download?url=VIDEO_URL

Get video information.

**Response:**
```json
{
  "title": "Video Title",
  "duration": "180",
  "thumbnail": "https://...",
  "author": "Channel Name",
  "views": "1000000",
  "formats": [...]
}
```

## ⚠️ Limitations

### Current Version (YouTube Only)

Due to serverless function limitations, the current web version only supports **YouTube**.

For other platforms (Vimeo, Twitter, TikTok, Instagram, etc.), please use the [CLI version](https://github.com/vib795/pull-vids).

### DRM Content

Cannot download from:
- Netflix, Prime Video, Disney+, Hulu
- HBO Max, Apple TV+, Paramount+
- Other paid streaming services with DRM

### File Size

Vercel free tier has these limits:
- Max function duration: 30 seconds
- Max response size: 4.5 MB

For large videos, use the CLI version.

## 🔮 Future Enhancements

- [ ] Support for more platforms (via proxy or worker)
- [ ] Batch downloads
- [ ] Playlist support
- [ ] Download history
- [ ] Dark mode toggle
- [ ] Video preview before download
- [ ] Format conversion
- [ ] Subtitle download

## 🤝 Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally
5. Submit a pull request

## 📄 License

MIT License - same as the main pull-vids project.

## 🔗 Links

- [Main Repository](https://github.com/vib795/pull-vids)
- [CLI Documentation](https://github.com/vib795/pull-vids#readme)
- [Homebrew Tap](https://github.com/vib795/homebrew-tap)
- [Report Issues](https://github.com/vib795/pull-vids/issues)

## 💡 Why Use CLI Instead?

The CLI version offers:
- Support for 1000+ platforms (not just YouTube)
- No file size limits
- Faster downloads
- Playlist support
- More format options
- Offline usage

**Install CLI:**
```bash
brew tap vib795/tap
brew install pull-vids
```

## 🙏 Acknowledgments

Built with:
- [Next.js](https://nextjs.org/)
- [ytdl-core](https://github.com/fent/node-ytdl-core)
- [Tailwind CSS](https://tailwindcss.com/)
- [Vercel](https://vercel.com/)

---

**Note:** This is a hobby project for educational purposes. Please respect copyright laws and platform Terms of Service.
