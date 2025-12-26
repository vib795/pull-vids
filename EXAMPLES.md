# pull-vids Examples

Practical examples for common use cases.

## Basic Usage

### Download a YouTube Video

```bash
pull-vids 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### Download from Vimeo

```bash
pull-vids 'https://vimeo.com/123456789'
```

### Download from Twitter/X

```bash
pull-vids 'https://twitter.com/user/status/1234567890'
```

### Download from TikTok

```bash
pull-vids 'https://www.tiktok.com/@username/video/1234567890'
```

### Download from Instagram

```bash
pull-vids 'https://www.instagram.com/p/ABC123/'
```

## Quality Selection

### Download Best Quality (Default)

```bash
pull-vids 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### Download 4K Video

```bash
pull-vids --quality 2160p 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### Download 1080p

```bash
pull-vids --quality 1080p 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### Download 720p (Medium Quality)

```bash
pull-vids --quality 720p 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### Download Low Quality (Faster)

```bash
pull-vids --quality low 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

## Audio Extraction

### Extract Audio as MP3

```bash
pull-vids --audio 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### Extract Audio as M4A

```bash
pull-vids --audio --format m4a 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### Download Music Video as Audio

```bash
pull-vids -a 'https://www.youtube.com/watch?v=MUSIC_VIDEO_ID'
```

## Playlists and Channels

### Download Entire Playlist

```bash
pull-vids --playlist 'https://www.youtube.com/playlist?list=PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
```

### Download Playlist as Audio

```bash
pull-vids --audio --playlist 'https://www.youtube.com/playlist?list=PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
```

### Download Channel's Videos

```bash
pull-vids --playlist 'https://www.youtube.com/@channelname/videos'
```

## Custom Output

### Save to Specific Directory

```bash
pull-vids --output ~/Videos 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### Save to Music Folder (Audio)

```bash
pull-vids --audio --output ~/Music 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### Save to Downloads

```bash
pull-vids -o ~/Downloads 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

## Format Options

### Download as MKV

```bash
pull-vids --format mkv 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### Download as WebM

```bash
pull-vids --format webm 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### Download Audio as WAV

```bash
pull-vids --audio --format wav 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

## Real-World Scenarios

### Download Tutorial Video in 1080p

```bash
pull-vids --quality 1080p --output ~/Tutorials 'https://www.youtube.com/watch?v=TUTORIAL_ID'
```

### Download Music Playlist as MP3

```bash
pull-vids --audio --format mp3 --playlist --output ~/Music/MyPlaylist \
  'https://www.youtube.com/playlist?list=PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
```

### Download Conference Talk

```bash
pull-vids --quality high --output ~/Conferences \
  'https://www.youtube.com/watch?v=CONFERENCE_TALK_ID'
```

### Download Stream VOD from Twitch

```bash
pull-vids 'https://www.twitch.tv/videos/1234567890'
```

### Download Cooking Recipe Video

```bash
pull-vids --quality 720p --output ~/Recipes \
  'https://www.youtube.com/watch?v=RECIPE_ID'
```

### Download Workout Video for Offline Viewing

```bash
pull-vids --quality 1080p --output ~/Workouts \
  'https://www.youtube.com/watch?v=WORKOUT_ID'
```

## Batch Downloads

### Download Multiple Videos (Script)

```bash
#!/bin/bash
# download-list.sh

videos=(
  "https://www.youtube.com/watch?v=VIDEO1"
  "https://www.youtube.com/watch?v=VIDEO2"
  "https://www.youtube.com/watch?v=VIDEO3"
)

for url in "${videos[@]}"; do
  pull-vids "$url"
done
```

### Download from File List

```bash
# Create urls.txt with one URL per line
cat urls.txt | while read url; do
  pull-vids "$url"
done
```

## Advanced Examples

### Download with Custom Quality and Format

```bash
pull-vids --quality 1080p --format mp4 --output ~/Videos \
  'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### Download Entire Playlist in Specific Quality

```bash
pull-vids --playlist --quality 720p --output ~/MyPlaylist \
  'https://www.youtube.com/playlist?list=PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
```

### Extract Audio from Playlist

```bash
pull-vids --audio --format mp3 --playlist --output ~/Music \
  'https://www.youtube.com/playlist?list=PLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
```

## Common Patterns

### Quick Download (Default Settings)

```bash
pull-vids 'VIDEO_URL'
```

### Audio Download

```bash
pull-vids -a 'VIDEO_URL'
```

### HD Download

```bash
pull-vids -q 1080p 'VIDEO_URL'
```

### Playlist Download

```bash
pull-vids -p 'PLAYLIST_URL'
```

### Custom Output

```bash
pull-vids -o ~/Videos 'VIDEO_URL'
```

## Troubleshooting Examples

### Check Version

```bash
pull-vids --version
```

### View Help

```bash
pull-vids --help
```

### Test with Short Video

```bash
# Use a short video to test
pull-vids 'https://www.youtube.com/watch?v=aqz-KE-bpKQ'  # Big Buck Bunny trailer
```

## Pro Tips

### 1. Use Quotes Around URLs

```bash
# Good
pull-vids 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'

# Bad (may fail with special characters)
pull-vids https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

### 2. Check Available Quality First

```bash
# yt-dlp can show available formats
yt-dlp -F 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```

### 3. Create Aliases for Common Tasks

```bash
# Add to ~/.bashrc or ~/.zshrc
alias yt-audio='pull-vids --audio'
alias yt-1080p='pull-vids --quality 1080p'
alias yt-playlist='pull-vids --playlist'

# Usage
yt-audio 'VIDEO_URL'
yt-1080p 'VIDEO_URL'
```

### 4. Organize by Category

```bash
# Music
pull-vids --audio --output ~/Music/YouTube 'MUSIC_URL'

# Tutorials
pull-vids --quality 1080p --output ~/Tutorials 'TUTORIAL_URL'

# Movies
pull-vids --quality best --output ~/Movies 'MOVIE_URL'
```

## Integration Examples

### Use in Shell Script

```bash
#!/bin/bash
# auto-download.sh

OUTPUT_DIR="$HOME/Downloads/Videos"
QUALITY="1080p"

if [ -z "$1" ]; then
  echo "Usage: $0 <video-url>"
  exit 1
fi

pull-vids --quality "$QUALITY" --output "$OUTPUT_DIR" "$1"
echo "Download complete! Saved to: $OUTPUT_DIR"
```

### Use with cron (Scheduled Downloads)

```bash
# Add to crontab: crontab -e
# Download daily at 2 AM
0 2 * * * /usr/local/bin/pull-vids --playlist 'PLAYLIST_URL' >> /var/log/pull-vids.log 2>&1
```

### Use in Makefile

```makefile
.PHONY: download-tutorials

download-tutorials:
	@echo "Downloading tutorial videos..."
	@pull-vids --quality 1080p --output ./tutorials 'TUTORIAL_PLAYLIST_URL'
```

---

For more examples and documentation, visit: https://github.com/vib795/pull-vids
