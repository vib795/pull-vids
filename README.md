# pull-vids

A free, open-source CLI tool for downloading YouTube videos and audio. Inspired by tools like Downie and PullTube, but completely free!

## Features

- Download YouTube videos in various qualities (360p to 4K)
- Download audio-only (MP3, M4A, etc.)
- Support for playlists
- Beautiful progress indicators
- Cross-platform (Windows, macOS, Linux)
- No ads, no paywalls, just downloads

## Installation

### Requirements

- Python 3.8 or higher
- ffmpeg (for audio conversion and video merging)

### Install ffmpeg

**macOS:**
```bash
brew install ffmpeg
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install ffmpeg
```

**Windows:**
Download from [ffmpeg.org](https://ffmpeg.org/download.html) or use:
```bash
winget install ffmpeg
```

### Install pull-vids

```bash
# Clone the repository
git clone https://github.com/vib795/pull-vids.git
cd pull-vids

# Install dependencies
pip install -r requirements.txt

# Install the package
pip install -e .
```

## Usage

### Basic Examples

**Download a video (best quality):**
```bash
pull-vids "https://www.youtube.com/watch?v=VIDEO_ID"
```

**Download audio only:**
```bash
pull-vids -a "https://www.youtube.com/watch?v=VIDEO_ID"
```

**Download in specific quality:**
```bash
pull-vids -q 720p "https://www.youtube.com/watch?v=VIDEO_ID"
```

**Download to a specific directory:**
```bash
pull-vids -o ~/Videos "https://www.youtube.com/watch?v=VIDEO_ID"
```

**Download entire playlist:**
```bash
pull-vids -p "https://www.youtube.com/playlist?list=PLAYLIST_ID"
```

**Download playlist as audio:**
```bash
pull-vids -a -p "https://www.youtube.com/playlist?list=PLAYLIST_ID"
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

Built with:
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) - The amazing YouTube downloading library
- [colorama](https://github.com/tartley/colorama) - Cross-platform colored terminal output

Inspired by Downie and PullTube, but free and open-source!

## Disclaimer

This tool is for personal use only. Please respect copyright laws and YouTube's Terms of Service. Only download videos you have the right to download.
