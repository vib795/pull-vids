# pull-vids

[![GitHub release](https://img.shields.io/github/v/release/vib795/pull-vids)](https://github.com/vib795/pull-vids/releases)
[![Homebrew](https://img.shields.io/badge/homebrew-vib795%2Ftap-orange)](https://github.com/vib795/homebrew-tap)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/vib795/pull-vids)](https://go.dev/)
[![GitHub stars](https://img.shields.io/github/stars/vib795/pull-vids?style=social)](https://github.com/vib795/pull-vids/stargazers)

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
- **Cookie authentication** for YouTube bot detection bypass
- **Automatic retry** with exponential backoff for rate limiting
- **Sleep intervals** between downloads to avoid rate limits
- Various quality options (360p to 4K)
- Audio-only extraction (MP3, M4A, etc.)
- **Transcripts** as plain text, SRT or VTT, without downloading the video
- Playlist and channel support
- Beautiful real-time progress bars
- Cross-platform (Windows, macOS, Linux)
- Single Go binary; the only runtime dependencies are yt-dlp and ffmpeg, which Homebrew and Chocolatey install for you
- No ads, no paywalls, just downloads

## Demo

![Demo](demo.gif)

## Installation

### Package Managers (Recommended)

**macOS (Homebrew):**
```bash
brew tap vib795/tap
brew install pull-vids
```

> **Upgrading from a version before 0.3.3?** Run this once, on each machine:
>
> ```bash
> brew uninstall pull-vids && brew install vib795/tap/pull-vids
> ```
>
> `brew upgrade` will report `pull-vids 64 already installed` and do nothing.
> See [Troubleshooting](#troubleshooting) for why.

**Windows (Chocolatey):**
```powershell
choco install pull-vids
```

**Any platform with Go 1.24+:**
```bash
go install github.com/vib795/pull-vids@latest
```

> **APT, AUR and Snap are not available yet.** Earlier versions of this README
> listed them before the packages existed, so `apt install pull-vids` and
> friends fail with "package was not found". Use Homebrew, Chocolatey,
> `go install`, or the install scripts below.

### Quick Install Script

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/vib795/pull-vids/develop/install.sh | bash
```

**Windows (PowerShell as Administrator):**
```powershell
irm https://raw.githubusercontent.com/vib795/pull-vids/develop/install.ps1 | iex
```

### Manual Installation

#### Option 1: Download Pre-built Binary

Download the latest release for your platform from [Releases](https://github.com/vib795/pull-vids/releases):

Releases ship as compressed archives, so download and extract before installing.

**macOS:**
```bash
# Apple Silicon (M1/M2/M3) - use darwin-amd64 on Intel
curl -L https://github.com/vib795/pull-vids/releases/latest/download/pull-vids-darwin-arm64.tar.gz -o pull-vids.tar.gz
tar -xzf pull-vids.tar.gz
chmod +x pull-vids-darwin-arm64
sudo mv pull-vids-darwin-arm64 /usr/local/bin/pull-vids
```

**Linux:**
```bash
# AMD64 - use linux-arm64 on ARM
curl -L https://github.com/vib795/pull-vids/releases/latest/download/pull-vids-linux-amd64.tar.gz -o pull-vids.tar.gz
tar -xzf pull-vids.tar.gz
chmod +x pull-vids-linux-amd64
sudo mv pull-vids-linux-amd64 /usr/local/bin/pull-vids
```

**Windows:**
1. Download [pull-vids-windows-amd64.zip](https://github.com/vib795/pull-vids/releases/latest/download/pull-vids-windows-amd64.zip)
2. Extract it, then rename `pull-vids-windows-amd64.exe` to `pull-vids.exe`
3. Move it to a directory in your PATH (e.g., `C:\Program Files\pull-vids\`)

#### Option 2: Build from Source

**Requirements:**
- Go 1.24 or higher (see `go.mod`)
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

### Transcripts and Captions

**Save just the transcript** (no video download, plain text by default):
```bash
pull-vids -t "https://www.youtube.com/watch?v=VIDEO_ID"
```

**Keep the timestamps** with `-f srt` or `-f vtt`, and choose a language with `--sub-langs`:
```bash
pull-vids -t -f srt --sub-langs es "https://www.youtube.com/watch?v=VIDEO_ID"
```

**Download a video with its captions** saved alongside it as `.srt`:
```bash
pull-vids --subs "https://www.youtube.com/watch?v=VIDEO_ID"
```

Captions uploaded by the creator are used when they exist; otherwise pull-vids
falls back to auto-generated ones. Plain-text transcripts drop the timestamps
and markup, and remove the repeated lines that YouTube's scrolling
auto-generated captions contain. `--sub-langs` defaults to `en`, takes
comma-separated codes or yt-dlp patterns such as `en.*`, and
`yt-dlp --list-subs <URL>` shows which languages a video offers.

### YouTube Authentication (Cookie Support)

**If you get a bot detection error from YouTube**, you'll need to authenticate using cookies from your browser.

#### Option 1: Export Cookies Manually (Recommended - Most Reliable)

1. **Install a browser extension to export cookies:**
   - **Chrome/Edge/Brave:** [Get cookies.txt LOCALLY](https://chrome.google.com/webstore/detail/get-cookiestxt-locally/cclelndahbckbenkjhflpdbgdldlbecc)
   - **Firefox:** [cookies.txt](https://addons.mozilla.org/en-US/firefox/addon/cookies-txt/)

2. **Visit YouTube and make sure you're logged in**

3. **Click the extension icon and export cookies** for `youtube.com`

4. **Save the file** (e.g., `youtube-cookies.txt`)

5. **Use with pull-vids:**
   ```bash
   pull-vids --cookies youtube-cookies.txt "https://www.youtube.com/watch?v=VIDEO_ID"
   ```

#### Option 2: Extract Cookies from Browser

**Firefox (usually works without issues):**
```bash
pull-vids --cookies-from-browser firefox "https://www.youtube.com/watch?v=VIDEO_ID"
```

**Chrome:**
```bash
pull-vids --cookies-from-browser chrome "https://www.youtube.com/watch?v=VIDEO_ID"
```

**Safari (macOS only):**
```bash
pull-vids --cookies-from-browser safari "https://www.youtube.com/watch?v=VIDEO_ID"
```

**⚠️ macOS Safari Users:** If you get a permission error, you need to grant Full Disk Access:
1. Open **System Settings** → **Privacy & Security** → **Full Disk Access**
2. Click the **+** button and add your terminal app (Terminal.app or iTerm2)
3. Restart your terminal
4. Try again

**⚠️ Chrome/Chromium Users:** If Chrome doesn't work, try exporting cookies manually (Option 1) or use Firefox.

**Supported browsers:** `firefox`, `chrome`, `safari`, `edge`, `chromium`, `brave`, `opera`, `vivaldi`

**Important:** Make sure you're logged into YouTube in the browser before extracting cookies!

### Command-Line Options

```
usage: pull-vids [-h] [-o OUTPUT] [-q QUALITY] [-a] [-p] [-f FORMAT]
                 [--cookies COOKIES] [--cookies-from-browser BROWSER]
                 [--sleep-interval SECONDS] [-N CONNECTIONS]
                 [--downloader BACKEND] [--http-chunk-size SIZE]
                 [-t] [--subs] [--sub-langs LANGS]
                 [-v] [--no-banner] url

positional arguments:
  url                   Video URL from any supported site (1000+ platforms)

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
                        Output format (mp4, mkv, mp3, m4a, etc.; with -t: txt, srt, vtt)
  -t, --transcript      Save the transcript only, without downloading the video
  --subs                Also save captions as .srt files alongside the video
  --sub-langs LANGS     Caption languages for -t and --subs (default: en)
  --cookies COOKIES     Path to cookies file (Netscape format)
  --cookies-from-browser BROWSER
                        Extract cookies from browser (chrome, firefox, safari, edge, etc.)
  --sleep-interval SECONDS
                        Sleep interval in seconds between downloads (avoids rate limiting)
  -N, --connections N   Parallel fragments for HLS/DASH formats, or connections
                        with aria2c; halved automatically on 403 (default: 8)
  --downloader BACKEND  Transfer backend: auto (native), native, or aria2c
                        (default: auto)
  --http-chunk-size SIZE
                        Chunk size for the native downloader (default: 10M;
                        keep it set, YouTube throttles unchunked requests)
  -v, --version         show program's version number and exit
  --no-banner           Don't show the banner
```

### Download Speed

YouTube throttles a download by the shape of its requests, not by how many
connections it uses. Asking for a large file in one request is held to roughly
playback speed, while asking for it in bounded pieces of 10 MB or less is
served at full speed, even over a single connection. On the same 230 MiB file,
one request ran at 1.8 MiB/s and 10 MB pieces at 42.7 MiB/s.

That is why `pull-vids` uses yt-dlp's native downloader with
`--http-chunk-size 10M` by default. aria2c is still available, but it measured
slower end to end at every size tried, partly because anything below its split
size goes out as one unranged request. Timings include yt-dlp's few seconds of
page extraction:

| Video                          | Size      | native (default) | `--downloader aria2c -N 8` |
|--------------------------------|-----------|------------------|----------------------------|
| 19-second clip                 | 0.5 MiB   | 5.1 s            | 22.3 s                     |
| 3.5-minute music video, 1080p  | 32.3 MiB  | 6.2 s            | 19.3 s                     |
| Same video, best quality (4K)  | 232.5 MiB | 13.3 s           | 34.8 s                     |

So there is nothing to install for speed. A few things still matter:

- **`-N` only helps fragmented streams.** HLS and DASH formats (Twitch, many
  news and sports sites) arrive as many small fragments, and `-N` downloads
  that many at once. YouTube's usual formats are single files, where the chunk
  size, not `-N`, sets the speed.
- **Don't clear `--http-chunk-size`.** An empty value turns chunking off and
  puts you back at playback speed on YouTube.
- **403 errors back off automatically.** If the server rejects a burst of
  parallel requests, `pull-vids` halves `-N` on each retry.
- **Storage caps everything.** Writing to a network share limits throughput
  to the share's write speed.

Run `./speedtest.sh` in a checkout to compare the backends on your own
connection.

### Quality Options

- `best` - Best available quality (default)
- `high` - Up to 1080p
- `medium` - Up to 720p
- `low` - Up to 480p
- `2160p`, `1440p`, `1080p`, `720p`, `480p`, `360p` - Specific resolutions

### Format Options

**Video formats:** by default the container follows the source (usually webm or mp4); `-f mp4`, `-f mkv` or `-f webm` picks it when streams are merged
**Audio formats:** mp3 (default), m4a, opus, wav
**Transcript formats (with `-t`):** txt (default), srt, vtt

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

# Download with cookies for YouTube bot detection
pull-vids --cookies-from-browser firefox "https://www.youtube.com/watch?v=VIDEO_ID"

# Download long playlist with rate limit protection
pull-vids --cookies-from-browser firefox --sleep-interval 5 -p "https://www.youtube.com/playlist?list=PLAYLIST_ID"
```

## Why pull-vids?

- **Free Forever** - No subscriptions, no trials, no limitations
- **Open Source** - Transparent code you can trust and modify
- **Privacy Focused** - No tracking, no data collection
- **Powerful** - Built on yt-dlp, supporting 1000+ video platforms
- **Simple** - Clean CLI interface, no bloat
- **Fast** - Requests streams the way YouTube serves at full speed
- **Portable** - One executable per platform; bring yt-dlp and ffmpeg

## Limitations

### DRM-Protected Content

**This tool CANNOT download DRM-protected content.** This includes:

- ❌ Amazon Prime Video
- ❌ Netflix
- ❌ Disney+
- ❌ Hulu
- ❌ HBO Max
- ❌ Apple TV+
- ❌ Other paid streaming services

**Why?** These services use DRM (Digital Rights Management) protection. yt-dlp does not and will not support bypassing DRM due to legal restrictions and ethical reasons.

### What You CAN Download

✅ **User-uploaded content:**
- YouTube (videos, playlists, channels, live streams)
- Vimeo (public and password-protected videos)
- Social media (Twitter, TikTok, Instagram, Facebook)
- Twitch (VODs and clips)
- Reddit (v.redd.it videos)
- And 1000+ other sites with publicly accessible content

✅ **Public platforms** where content creators share their work
✅ **Content you have legal rights to download**

**Note:** Just because you can download something doesn't mean you should. Always respect:
- Copyright laws
- Platform Terms of Service
- Content creator rights
- Fair use guidelines

## Troubleshooting

**`brew upgrade` says "pull-vids 64 already installed" and never updates:**

Run this once on the affected machine:

```bash
brew uninstall pull-vids && brew install vib795/tap/pull-vids
```

Formulas before 0.3.3 did not declare a version, so Homebrew inferred one from
the download filename — `pull-vids-darwin-arm64.tar.gz` yields `64`. Every
release looked like version `64`, so the binary was installed into
`Cellar/pull-vids/64` and each new release appeared to be already installed.

`brew upgrade` cannot repair this, because it compares the installed version
against the formula version and `64` sorts higher than `0.3.4`. Homebrew
concludes the installed copy is newer and declines. Only a fresh install
rebuilds the keg under the correct name.

Releases from 0.3.3 onward declare the version explicitly, so this is a
one-time cleanup. Machines installing fresh are unaffected and upgrade
normally.

**"ffmpeg not found" error:**
- Make sure ffmpeg is installed and in your PATH
- Try running `ffmpeg -version` to verify installation

**"No video formats found" error:**
- The video might be region-locked or unavailable
- Try a different quality setting
- Check if the URL is correct

**"Sign in to confirm you're not a bot" error:**
- See the [YouTube Authentication (Cookie Support)](#youtube-authentication-cookie-support) section above
- Use `--cookies-from-browser` or `--cookies` flag
- pull-vids will automatically retry up to 3 times with exponential backoff (30s, 60s, 120s)

**Rate limiting on long playlists:**
- Use `--sleep-interval 5` to add delays between downloads
- Recommended: 5-10 seconds for playlists with 30+ videos
- Automatic retry kicks in if rate limiting is detected despite the sleep interval

**Slow downloads:**
- Make sure `--http-chunk-size` hasn't been set to an empty value; unchunked
  requests are throttled to about playback speed
- If you passed `--downloader aria2c`, drop it; the native default is faster
  on YouTube
- Update yt-dlp (`brew upgrade yt-dlp`, `choco upgrade yt-dlp` or
  `pip install -U yt-dlp`); YouTube changes regularly break older builds
- Writing to a network share caps throughput at the share's write speed
- See [Download Speed](#download-speed) for measured numbers
- Use a different quality setting
- Some platforms have rate limits

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

This tool is for personal use only. Please respect copyright laws and each platform's Terms of Service. Only download videos you have the right to download.