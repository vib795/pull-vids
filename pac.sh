#!/bin/bash
# ~/bin/pull-and-convert (or anywhere in your PATH)

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default output directory (matches pull-vids default)
DOWNLOAD_DIR="${HOME}/Downloads/pull-vids"

# Pass all arguments to pull-vids
echo -e "${GREEN}📥 Downloading video...${NC}"
pull-vids "$@"

# Check if download succeeded
if [ $? -ne 0 ]; then
    echo -e "${YELLOW}Download failed or cancelled${NC}"
    exit 1
fi

# Find the most recently downloaded file
LATEST_FILE=$(ls -t "$DOWNLOAD_DIR"/*.{mp4,mkv,webm,avi,mov,flv} 2>/dev/null | head -1)

if [ -z "$LATEST_FILE" ]; then
    echo -e "${YELLOW}No video file found in $DOWNLOAD_DIR${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Downloaded: ${LATEST_FILE}${NC}"
echo ""

# Ask if user wants to convert
read -p "🔄 Convert this video to another format? (y/n): " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    # Show format options
    echo ""
    echo "Available formats: mp4, avi, mov, mkv, webm, flv"
    read -p "Target format (default: mp4): " FORMAT
    FORMAT=${FORMAT:-mp4}
    
    # Quality selection
    echo ""
    echo "Quality presets: high, medium, low"
    read -p "Quality (default: medium): " QUALITY
    QUALITY=${QUALITY:-medium}
    
    # Generate output filename
    BASENAME=$(basename "$LATEST_FILE" | sed 's/\.[^.]*$//')
    OUTPUT_FILE="${DOWNLOAD_DIR}/${BASENAME}_converted.${FORMAT}"
    
    echo -e "${GREEN}🎬 Converting to ${FORMAT}...${NC}"
    convert-vid convert "$LATEST_FILE" -f "$FORMAT" -q "$QUALITY" -o "$OUTPUT_FILE"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Converted: ${OUTPUT_FILE}${NC}"
        
        # Optionally ask to delete original
        read -p "🗑️  Delete original file? (y/n): " -n 1 -r
        echo ""
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm "$LATEST_FILE"
            echo "Original deleted."
        fi
    fi
fi
