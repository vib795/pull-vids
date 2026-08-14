#!/bin/bash
# Install script for pull-vids on Unix-like systems (Linux, macOS)

set -e

INSTALL_DIR="/usr/local/bin"
REPO="vib795/pull-vids"
BINARY_NAME="pull-vids"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}╔═══════════════════════════════════════╗${NC}"
echo -e "${CYAN}║     pull-vids Installation Script     ║${NC}"
echo -e "${CYAN}╚═══════════════════════════════════════╝${NC}"
echo ""

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}✗ Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

# Map OS to expected format
case "$OS" in
    darwin)
        PLATFORM="darwin"
        ;;
    linux)
        PLATFORM="linux"
        ;;
    *)
        echo -e "${RED}✗ Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

echo -e "${YELLOW}Detected platform: ${PLATFORM}-${ARCH}${NC}"
echo ""

# Check for curl or wget
if command -v curl &> /dev/null; then
    DOWNLOAD_CMD="curl -fsSL"
elif command -v wget &> /dev/null; then
    DOWNLOAD_CMD="wget -qO-"
else
    echo -e "${RED}✗ Neither curl nor wget found. Please install one of them.${NC}"
    exit 1
fi

# Get latest release
echo -e "${CYAN}Fetching latest release...${NC}"
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")

if [ -z "$LATEST_RELEASE" ]; then
    echo -e "${YELLOW}No release found. Building from source...${NC}"

    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}✗ Go is not installed. Please install Go first.${NC}"
        echo -e "${YELLOW}Visit: https://golang.org/dl/${NC}"
        exit 1
    fi

    # Clone and build
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    echo -e "${CYAN}Cloning repository...${NC}"
    git clone "https://github.com/${REPO}.git" . || {
        echo -e "${RED}✗ Failed to clone repository${NC}"
        exit 1
    }

    echo -e "${CYAN}Building binary...${NC}"
    go build -o "$BINARY_NAME" main.go || {
        echo -e "${RED}✗ Build failed${NC}"
        exit 1
    }

    BINARY_PATH="$TEMP_DIR/$BINARY_NAME"
else
    echo -e "${GREEN}✓ Latest release: ${LATEST_RELEASE}${NC}"

    # Download and unpack the release archive. Releases ship .tar.gz archives,
    # not bare binaries, so requesting the binary name directly returns 404.
    ASSET_NAME="${BINARY_NAME}-${PLATFORM}-${ARCH}.tar.gz"
    BINARY_URL="https://github.com/${REPO}/releases/download/${LATEST_RELEASE}/${ASSET_NAME}"
    TEMP_DIR=$(mktemp -d)
    ARCHIVE_PATH="$TEMP_DIR/$ASSET_NAME"
    BINARY_PATH="$TEMP_DIR/${BINARY_NAME}-${PLATFORM}-${ARCH}"

    echo -e "${CYAN}Downloading ${ASSET_NAME}...${NC}"
    if [ "$DOWNLOAD_CMD" = "curl -fsSL" ]; then
        curl -fsSL "$BINARY_URL" -o "$ARCHIVE_PATH" || {
            echo -e "${RED}✗ Download failed: ${BINARY_URL}${NC}"
            exit 1
        }
    else
        wget -qO "$ARCHIVE_PATH" "$BINARY_URL" || {
            echo -e "${RED}✗ Download failed: ${BINARY_URL}${NC}"
            exit 1
        }
    fi

    echo -e "${CYAN}Extracting...${NC}"
    tar -xzf "$ARCHIVE_PATH" -C "$TEMP_DIR" || {
        echo -e "${RED}✗ Extraction failed${NC}"
        exit 1
    }

    # Fall back to whatever executable the archive holds, in case the asset
    # naming changes in a future release.
    if [ ! -f "$BINARY_PATH" ]; then
        BINARY_PATH=$(find "$TEMP_DIR" -type f ! -name '*.tar.gz' | head -1)
        if [ -z "$BINARY_PATH" ]; then
            echo -e "${RED}✗ No binary found inside ${ASSET_NAME}${NC}"
            exit 1
        fi
    fi
fi

# Make binary executable
chmod +x "$BINARY_PATH"

# Install binary
echo -e "${CYAN}Installing to ${INSTALL_DIR}...${NC}"
if [ -w "$INSTALL_DIR" ]; then
    cp "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
else
    sudo cp "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
fi

# Cleanup
rm -rf "$TEMP_DIR"

echo ""
echo -e "${GREEN}╔═══════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  ✓ Installation completed!            ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}Run '${BINARY_NAME} --help' to get started!${NC}"
echo ""
echo -e "${CYAN}Note: You also need to install:${NC}"
echo -e "  - ffmpeg (for video processing)"
echo -e "  - yt-dlp (for downloading from 1000+ sites)"
echo ""
echo -e "${CYAN}Install dependencies:${NC}"
if [ "$PLATFORM" = "darwin" ]; then
    echo -e "  ${YELLOW}brew install ffmpeg yt-dlp${NC}"
elif [ "$PLATFORM" = "linux" ]; then
    echo -e "  ${YELLOW}sudo apt install ffmpeg && pip install yt-dlp${NC}"
fi
