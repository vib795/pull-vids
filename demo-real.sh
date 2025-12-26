#!/bin/bash
# Real Demo Script - Actual commands for recording
# Use this with asciinema or just run it to see pull-vids in action

# Exit on error
set -e

# Function to simulate typing and waiting
demo_pause() {
    sleep 2
}

# Function to show command before running
show_and_run() {
    echo "$ $@"
    sleep 1
    "$@"
    demo_pause
}

clear

# Banner
echo "╔════════════════════════════════════════════════════╗"
echo "║       pull-vids - Universal Video Downloader       ║"
echo "║          1000+ Websites • Free & Open Source       ║"
echo "╚════════════════════════════════════════════════════╝"
echo ""
demo_pause

# 1. Show version
echo "# Check version"
show_and_run pull-vids --version

# 2. Show help (first 20 lines)
echo "# Available options"
show_and_run pull-vids --help

echo ""
echo "# Ready to download from:"
echo "  • YouTube • Vimeo • Twitter • TikTok • Instagram"
echo "  • Facebook • Twitch • Reddit • And 1000+ more!"
echo ""
demo_pause

# 3. Example commands (shown but not executed)
echo "# Example: Download YouTube video"
echo "$ pull-vids 'https://www.youtube.com/watch?v=VIDEO_ID'"
echo ""
demo_pause

echo "# Example: Extract audio as MP3"
echo "$ pull-vids --audio 'https://www.youtube.com/watch?v=VIDEO_ID'"
echo ""
demo_pause

echo "# Example: Download in 1080p"
echo "$ pull-vids --quality 1080p 'https://www.youtube.com/watch?v=VIDEO_ID'"
echo ""
demo_pause

echo "# Example: Download entire playlist"
echo "$ pull-vids --playlist 'https://www.youtube.com/playlist?list=PLAYLIST_ID'"
echo ""
demo_pause

# Installation info
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Installation (macOS):"
echo "  $ brew tap vib795/tap"
echo "  $ brew install pull-vids"
echo ""
echo "Star on GitHub: https://github.com/vib795/pull-vids ⭐"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
