# pull-vids

A free, open-source CLI tool for downloading videos and audio from **1000+ websites**. Inspired by tools like Downie and PullTube, but completely free!

**Built with Go for blazingly fast performance! 🚀**

## Supported Platforms

Works with any site supported by yt-dlp, including:
- **YouTube** - Videos, playlists, channels
- **Vimeo** - Videos and channels
- **Twitter/X** - Videos and GIFs
- **TikTok** - Videos
- **Instagram** - Videos, Reels, Stories
- **Facebook** - Videos
- **Twitch** - VODs and clips
- **Reddit** - Videos from v.redd.it
- **Dailymotion** - Videos
- **And 1000+ more!** - [Full list](https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md)

## Features

- Download videos from 1000+ websites
- Various quality options (360p to 4K)
- Audio-only extraction (MP3, M4A, etc.)
- Playlist and channel support
- Beautiful real-time progress bars
- Cross-platform (Windows, macOS, Linux)
- Single binary - no dependencies to install
- Fast startup - < 1ms
- No ads, no paywalls, just downloads

## Installation

### Package Managers (Recommended)

**macOS (Homebrew):**
```bash
brew tap vib795/pull-vids
brew install pull-vids
```

**Windows (Chocolatey):**
```powershell
choco install pull-vids
```

**Windows (Scoop):**
```powershell
scoop bucket add pull-vids https://github.com/vib795/scoop-pull-vids
scoop install pull-vids
```

**Ubuntu/Debian (APT):**
```bash
sudo add-apt-repository ppa:vib795/pull-vids
sudo apt update
sudo apt install pull-vids
```

**Arch Linux (AUR):**
```bash
yay -S pull-vids
```

**Linux (Snap):**
```bash
sudo snap install pull-vids
```

### Quick Install Script

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/vib795/pull-vids/main/install.sh | bash
```

**Windows (PowerShell as Administrator):**
```powershell
irm https://raw.githubusercontent.com/vib795/pull-vids/main/install.ps1 | iex
```

### Manual Installation

#### Option 1: Download Pre-built Binary

Download the latest release for your platform from [Releases](https://github.com/vib795/pull-vids/releases):

**macOS:**
```bash
# Intel Mac
curl -L https://github.com/vib795/pull-vids/releases/latest/download/pull-vids-darwin-amd64 -o pull-vids
chmod +x pull-vids
sudo mv pull-vids /usr/local/bin/

# Apple Silicon (M1/M2/M3)
curl -L https://github.com/vib795/pull-vids/releases/latest/download/pull-vids-darwin-arm64 -o pull-vids
chmod +x pull-vids
sudo mv pull-vids /usr/local/bin/
```

**Linux:**
```bash
# AMD64
curl -L https://github.com/vib795/pull-vids/releases/latest/download/pull-vids-linux-amd64 -o pull-vids
chmod +x pull-vids
sudo mv pull-vids /usr/local/bin/

# ARM64
curl -L https://github.com/vib795/pull-vids/releases/latest/download/pull-vids-linux-arm64 -o pull-vids
chmod +x pull-vids
sudo mv pull-vids /usr/local/bin/
```

**Windows:**
1. Download [pull-vids-windows-amd64.exe](https://github.com/vib795/pull-vids/releases/latest/download/pull-vids-windows-amd64.exe)
2. Rename to `pull-vids.exe`
3. Move to a directory in your PATH (e.g., `C:\Program Files\pull-vids\`)

#### Option 2: Build from Source

**Requirements:**
- Go 1.18 or higher
- Git

**Build:**
```bash
# Clone the repository
git clone https://github.com/vib795/pull-vids.git
cd pull-vids

# Build for your platform
make build

# Or build for all platforms
make build-all

# Install to system
make install
```

### Install Dependencies

**pull-vids requires ffmpeg and yt-dlp to work:**

**macOS:**
```bash
brew install ffmpeg yt-dlp
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install ffmpeg
pip install yt-dlp
```

**Windows:**
```powershell
winget install ffmpeg
pip install yt-dlp
```

**Verify installation:**
```bash
ffmpeg -version
yt-dlp --version
pull-vids --version
```

## Usage

### Basic Examples

**Download from YouTube:**
```bash
pull-vids "https://www.youtube.com/watch?v=VIDEO_ID"
```

**Download from Vimeo:**
```bash
pull-vids "https://vimeo.com/123456789"
```

**Download from Twitter/X:**
```bash
pull-vids "https://twitter.com/user/status/123456"
```

**Download from TikTok:**
```bash
pull-vids "https://www.tiktok.com/@user/video/123456"
```

**Download from Instagram:**
```bash
pull-vids "https://www.instagram.com/p/ABC123/"
```

**Download from Twitch:**
```bash
pull-vids "https://www.twitch.tv/videos/123456789"
```

**Download audio only:**
```bash
pull-vids -a "https://www.youtube.com/watch?v=VIDEO_ID"
```

**Download in specific quality:**
```bash
pull-vids -q 720p "https://vimeo.com/123456789"
```

**Download to a specific directory:**
```bash
pull-vids -o ~/Videos "https://www.youtube.com/watch?v=VIDEO_ID"
```

**Download entire playlist:**
```bash
pull-vids -p "https://www.youtube.com/playlist?list=PLAYLIST_ID"
```

### Command-Line Options

```
usage: pull-vids [-h] [-o OUTPUT] [-q QUALITY] [-a] [-p] [-f FORMAT] [-v] [--no-banner] url

positional arguments:
  url                   YouTube video or playlist URL

options:
  -h, --help            show this help message and exit
  -o OUTPUT, --output OUTPUT
                        Output directory (default: ~/Downloads/pull-vids)
  -q QUALITY, --quality QUALITY
                        Video quality: best, high, medium, low, 2160p, 1440p, 1080p, 720p, 480p, 360p
                        (default: best)
  -a, --audio-only      Download audio only
  -p, --playlist        Download entire playlist
  -f FORMAT, --format FORMAT
                        Output format (mp4, mkv, mp3, m4a, etc.)
  -v, --version         show program's version number and exit
  --no-banner           Don't show the banner
```

### Quality Options

- `best` - Best available quality (default)
- `high` - Up to 1080p
- `medium` - Up to 720p
- `low` - Up to 480p
- `2160p`, `1440p`, `1080p`, `720p`, `480p`, `360p` - Specific resolutions

### Format Options

**Video formats:** mp4 (default), mkv, webm
**Audio formats:** mp3 (default), m4a, opus, wav

## Examples

```bash
# Download 4K video
pull-vids -q 2160p "https://www.youtube.com/watch?v=VIDEO_ID"

# Download as MP3
pull-vids -a -f mp3 "https://www.youtube.com/watch?v=VIDEO_ID"

# Download playlist to Music folder
pull-vids -a -p -o ~/Music "https://www.youtube.com/playlist?list=PLAYLIST_ID"

# Download in MKV format
pull-vids -f mkv "https://www.youtube.com/watch?v=VIDEO_ID"
```

## Why pull-vids?

- **Free Forever** - No subscriptions, no trials, no limitations
- **Open Source** - Transparent code you can trust and modify
- **Privacy Focused** - No tracking, no data collection
- **Powerful** - Built on yt-dlp, the best YouTube downloading library
- **Simple** - Clean CLI interface, no bloat
- **Fast** - Compiles to a single binary with zero startup time
- **Portable** - Single executable, no runtime dependencies

## Troubleshooting

**"ffmpeg not found" error:**
- Make sure ffmpeg is installed and in your PATH
- Try running `ffmpeg -version` to verify installation

**"No video formats found" error:**
- The video might be region-locked or unavailable
- Try a different quality setting
- Check if the URL is correct

**Slow downloads:**
- YouTube may be throttling your connection
- Try downloading at a different time
- Use a different quality setting

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

MIT License - feel free to use this tool however you want!

## Credits

**Core:**
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) - The amazing YouTube downloading library

**Go libraries:**
- [fatih/color](https://github.com/fatih/color) - Colored terminal output
- [schollz/progressbar](https://github.com/schollz/progressbar) - Beautiful progress bars

Inspired by Downie and PullTube, but free and open-source!

## Disclaimer

This tool is for personal use only. Please respect copyright laws and YouTube's Terms of Service. Only download videos you have the right to download.
