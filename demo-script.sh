#!/bin/bash
# pull-vids Demo Script
# This script demonstrates the key features of pull-vids

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print section headers
print_header() {
    echo ""
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

# Function to run command with description
run_demo() {
    echo -e "${YELLOW}$ $1${NC}"
    sleep 1
    eval "$1"
    sleep 2
}

# Clear screen
clear

# Title
print_header "pull-vids - Universal Video Downloader"
echo -e "${GREEN}Download videos from 1000+ websites - Free & Open Source!${NC}"
sleep 2

# Demo 1: Show version
print_header "1. Check Version"
run_demo "pull-vids --version"

# Demo 2: Show help
print_header "2. Available Options"
run_demo "pull-vids --help"

# Demo 3: Download from YouTube (example)
print_header "3. Download YouTube Video"
echo "Example command:"
echo -e "${YELLOW}$ pull-vids 'https://www.youtube.com/watch?v=K14hxdekrzE'${NC}"
echo ""
echo "This will download the video in best available quality"
sleep 3

# Demo 4: Audio extraction
print_header "4. Extract Audio Only"
echo "Example command:"
echo -e "${YELLOW}$ pull-vids --audio 'https://www.youtube.com/watch?v=K14hxdekrzE'${NC}"
echo ""
echo "Downloads audio as MP3"
sleep 3

# Demo 5: Specific quality
print_header "5. Choose Specific Quality"
echo "Example command:"
echo -e "${YELLOW}$ pull-vids --quality 720p 'https://www.youtube.com/watch?v=K14hxdekrzE'${NC}"
echo ""
echo "Available qualities: best, high, medium, low, 2160p, 1440p, 1080p, 720p, 480p, 360p"
sleep 3

# Demo 6: Custom output directory
print_header "6. Custom Output Directory"
echo "Example command:"
echo -e "${YELLOW}$ pull-vids --output ~/Videos 'https://www.youtube.com/watch?v=K14hxdekrzE'${NC}"
echo ""
echo "Save to a specific folder"
sleep 3

# # Demo 7: Playlist download
# print_header "7. Download Entire Playlist"
# echo "Example command:"
# echo -e "${YELLOW}$ pull-vids --playlist 'https://www.youtube.com/playlist?list=PLAYLIST_ID'${NC}"
# echo ""
# echo "Downloads all videos in the playlist"
# sleep 3

# # Demo 8: Supported platforms
# print_header "8. Supported Platforms (1000+)"
# echo "✅ YouTube      - Videos, playlists, channels"
# echo "✅ Vimeo        - Videos and channels"
# echo "✅ Twitter/X    - Videos and GIFs"
# echo "✅ TikTok       - Videos"
# echo "✅ Instagram    - Videos, Reels, Stories"
# echo "✅ Facebook     - Videos"
# echo "✅ Twitch       - VODs and clips"
# echo "✅ Reddit       - v.redd.it videos"
# echo "✅ And 1000+ more!"
# sleep 4

# Installation reminder
print_header "Installation"
echo -e "${GREEN}macOS (Homebrew):${NC}"
echo -e "${YELLOW}$ brew tap vib795/tap${NC}"
echo -e "${YELLOW}$ brew install pull-vids${NC}"
echo ""
echo -e "${GREEN}Or download binaries from:${NC}"
echo "https://github.com/vib795/pull-vids/releases"
sleep 3

# Final message
print_header "Thank You!"
echo -e "${GREEN}⭐ Star on GitHub: https://github.com/vib795/pull-vids${NC}"
echo -e "${BLUE}📖 Full documentation in README.md${NC}"
echo -e "${CYAN}🚀 Built with Go for blazingly fast performance${NC}"
echo ""

exit 0
