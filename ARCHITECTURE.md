# Complete Guide to pull-vids: Code, Setup, and Package Management

## Table of Contents
1. [Code Architecture](#code-architecture)
2. [How the Application Works](#how-it-works)
3. [Setup Process](#setup-process)
4. [Package Management Explained](#package-management)
5. [Publishing Guide](#publishing-guide)

---

## 1. Code Architecture

### Project Structure
```
pull-vids/
├── main.go                      # Main application code (Go)
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── Makefile                     # Build automation
├── README.md                    # User documentation
├── LICENSE                      # MIT license
├── .gitignore                   # Git ignore rules
├── install.sh                   # Unix install script
├── install.ps1                  # Windows install script
├── PACKAGE_MANAGERS.md          # Publishing guide
├── Formula/
│   └── pull-vids.rb            # Homebrew formula
├── packaging/
│   ├── chocolatey/             # Chocolatey package
│   │   ├── pull-vids.nuspec
│   │   └── tools/
│   │       ├── chocolateyinstall.ps1
│   │       └── chocolateyuninstall.ps1
│   └── deb/                    # Debian package
│       └── DEBIAN/
│           ├── control
│           └── postinst
└── .github/
    └── workflows/
        └── release.yml         # Automated releases
```

### main.go Breakdown

#### **1. Imports and Dependencies**
```go
package main

import (
    "bufio"           // Reading lines from command output
    "encoding/json"   // Parsing yt-dlp JSON output
    "flag"            // Command-line argument parsing
    "fmt"             // Formatted I/O
    "os"              // Operating system functions
    "os/exec"         // Running external commands (yt-dlp)
    "path/filepath"   // Cross-platform path handling
    "strings"         // String manipulation
    "time"            // Timing downloads

    "github.com/fatih/color"                // Colored terminal output
    "github.com/schollz/progressbar/v3"     // Progress bars
)
```

**Why these libraries?**
- `fatih/color`: Makes terminal output colorful and user-friendly
- `progressbar`: Shows real-time download progress
- Standard library: Go's built-in packages handle everything else

#### **2. Configuration Struct**
```go
type Config struct {
    URL        string  // YouTube URL to download
    Output     string  // Where to save files
    Quality    string  // Video quality (720p, 1080p, etc.)
    AudioOnly  bool    // Download only audio?
    Playlist   bool    // Allow playlist downloads?
    Format     string  // Output format (mp4, mp3, etc.)
    NoBanner   bool    // Hide the banner?
    ShowVersion bool   // Show version and exit?
}
```

This struct holds all user preferences from command-line arguments.

#### **3. Key Functions Explained**

##### **cleanURL(url string)**
```go
func cleanURL(url string) string {
    // Removes backslash escapes from shell
    // Example: "watch\?v\=ABC" → "watch?v=ABC"
    replacements := map[string]string{
        `\?`: `?`,
        `\=`: `=`,
        `\&`: `&`,
    }
    // ... applies replacements
}
```

**Why needed?** When you type URLs in terminal without quotes, shells add backslashes:
```bash
# Without quotes, shell escapes special characters
./pull-vids https://youtube.com/watch\?v\=ABC

# cleanURL() removes those backslashes before passing to yt-dlp
```

##### **getFormatString(quality, audioOnly)**
```go
func getFormatString(quality string, audioOnly bool) string {
    if audioOnly {
        return "bestaudio/best"
    }

    qualityMap := map[string]string{
        "best":   "bestvideo+bestaudio/best",
        "1080p":  "bestvideo[height<=1080]+bestaudio/best[height<=1080]",
        "720p":   "bestvideo[height<=720]+bestaudio/best[height<=720]",
        // ... more quality options
    }

    return qualityMap[quality]
}
```

**What it does:** Translates simple quality names (like "720p") into yt-dlp format strings.
- `bestvideo[height<=720]`: Best video up to 720p
- `+bestaudio`: Merge with best audio
- `/best`: Fallback if specific format unavailable

##### **downloadVideo() - The Main Function**
```go
func downloadVideo(config *Config) error {
    // 1. Create output directory
    outputPath := os.ExpandEnv(config.Output)
    os.MkdirAll(outputPath, 0755)

    // 2. Build yt-dlp command
    args := []string{
        "--newline",                              // Progress on new lines
        "--progress",                             // Show progress
        "-f", getFormatString(config.Quality),    // Quality setting
        "-o", outputTemplate,                     // Where to save
    }

    // 3. Add audio extraction if needed
    if config.AudioOnly {
        args = append(args, "-x", "--audio-format", "mp3")
    }

    // 4. Execute yt-dlp
    cmd := exec.Command("yt-dlp", args...)

    // 5. Create pipes to read output
    stdout, _ := cmd.StdoutPipe()
    stderr, _ := cmd.StderrPipe()

    // 6. Start command
    cmd.Start()

    // 7. Parse output for progress
    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        line := scanner.Text()
        // Extract percentage, speed, ETA
        // Update progress bar
    }

    // 8. Wait for completion
    cmd.Wait()
}
```

**How it works:**
1. **Validates** output directory exists
2. **Builds** yt-dlp command with user options
3. **Executes** yt-dlp as a subprocess
4. **Reads** yt-dlp's output in real-time
5. **Parses** download percentage, speed, ETA
6. **Updates** progress bar continuously
7. **Reports** completion or errors

##### **Progress Bar Implementation**
```go
// Create a 100% scale progress bar
bar := progressbar.NewOptions(100,
    progressbar.OptionSetDescription("Downloading"),
    progressbar.OptionSetWidth(50),
    progressbar.OptionSetTheme(progressbar.Theme{
        Saucer:        "█",  // Filled part
        SaucerPadding: " ",  // Empty part
    }),
)

// Parse yt-dlp output
// "[download]  45.3% of 12.34MiB at 1.23MiB/s ETA 00:10"
if strings.Contains(line, "%") {
    percentStr := extractPercentage(line)  // "45.3%"
    percent := parseFloat(percentStr)       // 45.3
    bar.Set(int(percent))                   // Update to 45%
}
```

---

## 2. How It Works (Step by Step)

### User runs command:
```bash
pull-vids -q 720p "https://youtube.com/watch?v=ABC"
```

### What happens internally:

1. **Argument Parsing**
   ```go
   flag.StringVar(&config.Quality, "q", "best", "Video quality")
   flag.Parse()
   // config.Quality = "720p"
   // config.URL = "https://youtube.com/watch?v=ABC"
   ```

2. **URL Cleaning**
   ```go
   cleanURL := cleanURL(args[0])
   // Removes any shell escapes
   ```

3. **Validation**
   ```go
   if !strings.Contains(url, "youtube.com") {
       return errors.New("Invalid YouTube URL")
   }
   ```

4. **Build yt-dlp Command**
   ```go
   cmd := exec.Command("yt-dlp",
       "-f", "bestvideo[height<=720]+bestaudio",
       "-o", "~/Downloads/pull-vids/%(title)s.%(ext)s",
       "https://youtube.com/watch?v=ABC"
   )
   ```

5. **Execute & Monitor**
   - Start yt-dlp
   - Read its output line by line
   - Parse: `[download] 45.3% of 12.34MiB at 1.23MiB/s ETA 00:10`
   - Extract: 45.3%, 1.23MiB/s, 00:10
   - Update progress bar

6. **Complete**
   - yt-dlp finishes downloading
   - yt-dlp converts/merges video+audio
   - File saved to disk
   - Show success message

---

## 3. Setup Process

### For Developers (Building from Source)

#### **Step 1: Install Prerequisites**

**macOS:**
```bash
# Install Homebrew if not installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install dependencies
brew install go ffmpeg yt-dlp
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install golang-go ffmpeg
pip install yt-dlp
```

**Windows:**
```powershell
# Using winget
winget install GoLang.Go
winget install ffmpeg
pip install yt-dlp
```

#### **Step 2: Clone Repository**
```bash
git clone https://github.com/vib795/pull-vids.git
cd pull-vids
```

#### **Step 3: Understanding Go Modules**

**`go.mod` file:**
```go
module github.com/vib795/pull-vids

go 1.21

require (
    github.com/fatih/color v1.18.0
    github.com/schollz/progressbar/v3 v3.18.0
)
```

This declares:
- **Module name**: Used for imports
- **Go version**: Minimum required version
- **Dependencies**: External libraries needed

**Install dependencies:**
```bash
go mod download  # Download dependencies
go mod tidy      # Clean up unused deps
```

#### **Step 4: Build**

**Option A: Using Make**
```bash
make build
# Runs: go build -o pull-vids main.go
```

**Option B: Using Go directly**
```bash
go build -o pull-vids main.go
```

**What happens during build:**
1. Go compiler reads `main.go`
2. Resolves imports from `go.mod`
3. Downloads dependencies if needed
4. Compiles everything into single binary
5. Links external libraries statically
6. Produces `pull-vids` executable

**Build flags explained:**
```bash
go build \
  -ldflags "-s -w -X main.version=v0.2.0" \
  -o pull-vids \
  main.go
```

- `-ldflags`: Linker flags
  - `-s`: Strip debug symbols (smaller binary)
  - `-w`: Strip DWARF debug info (smaller binary)
  - `-X main.version=v0.2.0`: Inject version at build time
- `-o pull-vids`: Output filename

#### **Step 5: Test**
```bash
./pull-vids --version
# Output: pull-vids 0.2.0

./pull-vids --help
# Shows all options
```

#### **Step 6: Install System-Wide**
```bash
# macOS/Linux
sudo make install
# Copies to /usr/local/bin/pull-vids

# Windows
# Manually copy to C:\Program Files\pull-vids\
# Add to PATH
```

---

## 4. Package Management Explained

Package managers automate installation for users. Instead of:
```bash
git clone → go build → sudo mv → configure PATH
```

Users just type:
```bash
brew install pull-vids
```

### Homebrew (macOS/Linux)

#### **How Homebrew Works**

1. **Formula** = Recipe for installing software
2. **Tap** = Repository of formulas
3. **Cellar** = Where Homebrew stores installed software

#### **Our Formula: `Formula/pull-vids.rb`**

```ruby
class PullVids < Formula
  desc "Free CLI tool for downloading YouTube videos"
  homepage "https://github.com/vib795/pull-vids"
  version "0.2.0"
  license "MIT"

  # Platform-specific downloads
  on_macos do
    if Hardware::CPU.arm?
      # Apple Silicon (M1/M2/M3)
      url "https://github.com/.../pull-vids-darwin-arm64.tar.gz"
      sha256 "abc123..."
    else
      # Intel Mac
      url "https://github.com/.../pull-vids-darwin-amd64.tar.gz"
      sha256 "def456..."
    end
  end

  # Dependencies
  depends_on "ffmpeg"
  depends_on "yt-dlp"

  def install
    # Extract binary from archive
    # Move to Homebrew's bin directory
    bin.install "pull-vids-darwin-arm64" => "pull-vids"
  end

  test do
    # Verification test
    assert_match "pull-vids", shell_output("#{bin}/pull-vids --version")
  end
end
```

#### **What each part does:**

1. **Platform detection**: Downloads correct binary for user's CPU
2. **Dependencies**: Automatically installs ffmpeg and yt-dlp
3. **Installation**: Extracts and symlinks binary
4. **Test**: Verifies installation works

#### **Publishing to Homebrew**

**Option 1: Personal Tap (Easy)**
```bash
# Create your own tap repository
brew tap-new vib795/pull-vids

# Copy formula
cp Formula/pull-vids.rb \
  $(brew --repository)/Library/Taps/vib795/homebrew-pull-vids/Formula/

# Users install with:
brew tap vib795/pull-vids
brew install pull-vids
```

**Option 2: Homebrew Core (Official, but harder)**
- Submit PR to https://github.com/Homebrew/homebrew-core
- Stricter requirements
- More users discover it

### Chocolatey (Windows)

#### **How Chocolatey Works**

1. **NuSpec** = Package metadata (like Homebrew formula)
2. **Install Script** = PowerShell script to download & install
3. **Repository** = chocolatey.org hosts packages

#### **Our Package: `packaging/chocolatey/`**

**`pull-vids.nuspec` (Metadata):**
```xml
<?xml version="1.0"?>
<package>
  <metadata>
    <id>pull-vids</id>
    <version>0.2.0</version>
    <authors>pull-vids contributors</authors>
    <description>
      Free CLI tool for downloading YouTube videos
    </description>
    <dependencies>
      <dependency id="ffmpeg" version="4.4.0" />
    </dependencies>
  </metadata>
</package>
```

**`tools/chocolateyinstall.ps1` (Installation):**
```powershell
$packageArgs = @{
  packageName   = 'pull-vids'
  url64bit      = 'https://github.com/.../pull-vids-windows-amd64.zip'
  checksum64    = 'abc123...'
  checksumType64= 'sha256'
  unzipLocation = $toolsDir
}

# Download and extract
Install-ChocolateyZipPackage @packageArgs

# Rename binary
Rename-Item "pull-vids-windows-amd64.exe" "pull-vids.exe"
```

#### **Publishing to Chocolatey**

```powershell
# 1. Create account at chocolatey.org
# 2. Get API key from account settings

# 3. Build package
cd packaging/chocolatey
choco pack

# 4. Test locally
choco install pull-vids -s . -y

# 5. Publish
choco push pull-vids.0.2.0.nupkg \
  --source https://push.chocolatey.org/ \
  --api-key YOUR_API_KEY
```

### APT/Debian Packages

#### **How .deb Packages Work**

1. **Control File** = Package metadata
2. **Binary** = The actual program
3. **Scripts** = Run during install/uninstall

#### **Creating .deb Package**

```bash
# 1. Make build command
make deb

# What it does:
# Creates structure:
pull-vids_0.2.0_amd64/
├── usr/local/bin/
│   └── pull-vids           # Binary
└── DEBIAN/
    ├── control             # Metadata
    └── postinst            # Post-install script

# 2. Build package
dpkg-deb --build pull-vids_0.2.0_amd64

# 3. Result
pull-vids_0.2.0_amd64.deb
```

**`DEBIAN/control`:**
```
Package: pull-vids
Version: 0.2.0
Architecture: amd64
Depends: ffmpeg          # APT will install ffmpeg first
Description: Free CLI tool for downloading YouTube videos
```

#### **Publishing to Ubuntu PPA**

```bash
# 1. Create Launchpad account
# 2. Set up GPG keys for signing
# 3. Create PPA

# 4. Upload
dput ppa:vib795/pull-vids pull-vids_0.2.0_amd64.changes

# Users install with:
sudo add-apt-repository ppa:vib795/pull-vids
sudo apt update
sudo apt install pull-vids
```

---

## 5. Publishing Guide (Complete Workflow)

### Step 1: Prepare Release

```bash
# 1. Update version
VERSION="0.2.0"

# Update in:
# - main.go (const version)
# - Formula/pull-vids.rb
# - packaging/chocolatey/pull-vids.nuspec
# - packaging/deb/DEBIAN/control

# 2. Commit changes
git add .
git commit -m "Bump version to v${VERSION}"
git push

# 3. Create git tag
git tag "v${VERSION}"
git push origin "v${VERSION}"
```

### Step 2: Build All Binaries

```bash
make build-all
# Creates:
# - pull-vids-linux-amd64
# - pull-vids-linux-arm64
# - pull-vids-darwin-amd64
# - pull-vids-darwin-arm64
# - pull-vids-windows-amd64.exe
```

### Step 3: Create Release Archives

```bash
make release
# Creates compressed archives in dist/releases/:
# - pull-vids-linux-amd64.tar.gz
# - pull-vids-darwin-arm64.tar.gz
# - pull-vids-windows-amd64.zip
# etc.
```

### Step 4: Generate Checksums

```bash
make checksums
# Creates dist/releases/SHA256SUMS with:
# abc123... pull-vids-linux-amd64.tar.gz
# def456... pull-vids-darwin-arm64.tar.gz
# etc.
```

### Step 5: Create Debian Package

```bash
make deb
# Creates dist/packages/pull-vids_0.2.0_amd64.deb
```

### Step 6: GitHub Release (Automated)

When you push a tag, `.github/workflows/release.yml` automatically:

```yaml
# Triggered by tag push (v0.2.0)
on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    steps:
      - Build all binaries
      - Create archives
      - Generate checksums
      - Upload to GitHub Releases
```

Users can then download from:
`https://github.com/vib795/pull-vids/releases/tag/v0.2.0`

### Step 7: Update Homebrew

```bash
# Get checksums
shasum -a 256 dist/releases/pull-vids-darwin-*.tar.gz

# Update Formula/pull-vids.rb:
# - version = "0.2.0"
# - url = "...v0.2.0/pull-vids-darwin-arm64.tar.gz"
# - sha256 = "new_checksum"

# If you have a tap repository:
cd $(brew --repository)/Library/Taps/vib795/homebrew-pull-vids
git add Formula/pull-vids.rb
git commit -m "Update to v0.2.0"
git push
```

### Step 8: Publish to Chocolatey

```powershell
cd packaging/chocolatey

# Update checksums in tools/chocolateyinstall.ps1
# (Get-FileHash -Path pull-vids-windows-amd64.zip).Hash

choco pack
choco push pull-vids.0.2.0.nupkg `
  --source https://push.chocolatey.org/ `
  --api-key YOUR_API_KEY
```

### Step 9: Publish to APT/PPA

```bash
# Upload .deb to PPA
dput ppa:vib795/pull-vids dist/packages/pull-vids_0.2.0_amd64.deb
```

---

## Key Concepts Summary

### Cross-Compilation
```bash
# Build for different platforms from one machine
GOOS=linux GOARCH=amd64 go build -o pull-vids-linux
GOOS=darwin GOARCH=arm64 go build -o pull-vids-mac-m1
GOOS=windows GOARCH=amd64 go build -o pull-vids.exe
```

### Why Single Binary?
- Go compiles everything (including dependencies) into one file
- No need to install Node.js, Python, or other runtimes
- Just download and run

### Package Manager Benefits
| Without | With Package Manager |
|---------|---------------------|
| Manual download | `brew install pull-vids` |
| Manual updates | `brew upgrade pull-vids` |
| Manual dependencies | Auto-installs ffmpeg, yt-dlp |
| Complex uninstall | `brew uninstall pull-vids` |

### Versioning
- Git tags: `v0.2.0`
- Semantic versioning: MAJOR.MINOR.PATCH
  - MAJOR: Breaking changes
  - MINOR: New features
  - PATCH: Bug fixes

---

## Common Questions

**Q: Why use yt-dlp instead of implementing YouTube downloading ourselves?**
A: YouTube's API and download mechanisms are complex and constantly changing. yt-dlp is maintained by a large community and handles all edge cases.

**Q: Why Go instead of Python/Node.js?**
A:
- Single binary (no runtime needed)
- Fast startup (< 1ms vs 200-500ms for Python)
- Cross-compilation (build for all platforms from one machine)
- Better performance

**Q: How does progress bar work without yt-dlp API?**
A: We run yt-dlp as subprocess and parse its text output in real-time.

**Q: Do users need to install Go?**
A: No! The compiled binary has everything included. Only developers need Go.

**Q: What if GitHub Releases goes down?**
A: Package managers cache binaries. Plus we can host on other CDNs.

**Q: How to handle breaking changes in yt-dlp?**
A: Our code just passes arguments to yt-dlp. As long as yt-dlp's CLI doesn't change drastically, we're fine. If it does, we update our argument building logic.

---

This is the complete system! You now have:
1. ✅ Fast Go CLI tool
2. ✅ Cross-platform builds
3. ✅ Installation scripts
4. ✅ Package manager configurations
5. ✅ Automated release pipeline

Users can install from 6+ different package managers, and you can publish updates with a single `git tag` command!
