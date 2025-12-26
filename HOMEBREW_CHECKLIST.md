# Homebrew Core Contribution Guide - Complete Workflow

This guide walks you through **every step** needed to prepare, test, and submit your formula to Homebrew core.

---

## 🚀 Complete Workflow (Start to Finish)

Follow these steps **in order** on your Mac:

### Step 1: Pull Latest Changes

```bash
# Navigate to your project
cd ~/i/pull-vids  # or wherever your project is located

# Pull the latest changes with updated formula
git fetch origin
git checkout claude/universal-video-support-xSgrP
git pull origin claude/universal-video-support-xSgrP
```

---

### Step 2: Build Binaries on Mac

```bash
# Make sure Go is installed
go version  # Should show go1.21 or higher
# If not: brew install go

# Clean previous builds
rm -rf dist/
mkdir -p dist

# Build all platform binaries
make build-all

# Verify binaries were created
ls -lh dist/
# You should see:
# - pull-vids-darwin-amd64
# - pull-vids-darwin-arm64
# - pull-vids-linux-amd64
# - pull-vids-linux-arm64
# - pull-vids-windows-amd64.exe
```

---

### Step 3: Create Release Archives

```bash
# Create tar.gz archives for Homebrew
cd dist
tar -czf pull-vids-darwin-arm64.tar.gz pull-vids-darwin-arm64
tar -czf pull-vids-darwin-amd64.tar.gz pull-vids-darwin-amd64
tar -czf pull-vids-linux-arm64.tar.gz pull-vids-linux-arm64
tar -czf pull-vids-linux-amd64.tar.gz pull-vids-linux-amd64

# Verify archives were created
ls -lh *.tar.gz *.exe
# You should see all 5 files ready for release

cd ..  # Return to project root
```

---

### Step 4: Calculate SHA256 Checksums

```bash
# Calculate checksums for the archives
cd dist
shasum -a 256 *.tar.gz

# Copy the output - you'll need these checksums!
# Example output:
# a6d7c093bdd5cf387b5651f36dcdb92fec924f4ba6f2307383c0d15834bd5165  pull-vids-darwin-arm64.tar.gz
# 596ac58187b127a549d51bf079c4404837b39550b0fd9a816f18d2bbf3eb6af5  pull-vids-darwin-amd64.tar.gz
# f171eadb1884da3b4687a424c01fef56d2ea4394458319221b5c1d0e6e5bac23  pull-vids-linux-arm64.tar.gz
# 2146527a437c5dc56ae37ef004a302000a1290eb2affa5611bc86a9b1ac92839  pull-vids-linux-amd64.tar.gz

cd ..
```

---

### Step 5: Update Formula with Correct Checksums

**IMPORTANT:** Only do this if the checksums in `Formula/pull-vids.rb` don't match the ones you just calculated.

Edit `Formula/pull-vids.rb` and update the `sha256` values to match your calculated checksums:

```ruby
on_macos do
  on_arm do
    url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-darwin-arm64.tar.gz"
    sha256 "YOUR_DARWIN_ARM64_CHECKSUM_HERE"
  end
  on_intel do
    url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-darwin-amd64.tar.gz"
    sha256 "YOUR_DARWIN_AMD64_CHECKSUM_HERE"
  end
end

on_linux do
  on_arm do
    url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-linux-arm64.tar.gz"
    sha256 "YOUR_LINUX_ARM64_CHECKSUM_HERE"
  end
  on_intel do
    url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-linux-amd64.tar.gz"
    sha256 "YOUR_LINUX_AMD64_CHECKSUM_HERE"
  end
end
```

If you made changes, commit them:
```bash
git add Formula/pull-vids.rb
git commit -m "Update SHA256 checksums for v0.2.0 release"
git push origin claude/universal-video-support-xSgrP
```

---

### Step 6: Clean Up Old Releases (If Needed)

```bash
# Check if v0.2.0 or v0.2.1 already exist
gh release list

# If v0.2.0 exists with wrong files, delete it
gh release delete v0.2.0 --yes

# If v0.2.1 exists, delete it (we want v0.2.0)
gh release delete v0.2.1 --yes
```

---

### Step 7: Create GitHub Release v0.2.0

```bash
# Create the official v0.2.0 release
gh release create v0.2.0 \
  dist/pull-vids-darwin-arm64.tar.gz \
  dist/pull-vids-darwin-amd64.tar.gz \
  dist/pull-vids-linux-arm64.tar.gz \
  dist/pull-vids-linux-amd64.tar.gz \
  dist/pull-vids-windows-amd64.exe \
  --title "v0.2.0 - Universal Video Support" \
  --notes "Universal video downloader supporting 1000+ websites including YouTube, Vimeo, Twitter, TikTok, Instagram, Facebook, Twitch, Reddit, and more!

## Installation

### macOS (Homebrew)
\`\`\`bash
brew tap vib795/tap
brew install vib795/tap/pull-vids
\`\`\`

### Linux
Download the appropriate binary for your architecture.

### Windows
Download pull-vids-windows-amd64.exe

## Features
- ✅ Download from 1000+ websites
- ✅ Multiple quality options
- ✅ Audio extraction (MP3)
- ✅ Playlist support
- ✅ Progress bars

## Requirements
- ffmpeg
- yt-dlp"

# Verify the release was created
gh release view v0.2.0
```

---

### Step 8: Set Up Local Homebrew Tap

```bash
# Create a local tap for testing
brew tap-new vib795/tap

# Find the tap directory
TAP_DIR=$(brew --repository)/Library/Taps/vib795/homebrew-tap

# Copy your formula to the tap
cp Formula/pull-vids.rb $TAP_DIR/Formula/pull-vids.rb

# Verify it's there
ls -la $TAP_DIR/Formula/pull-vids.rb
```

---

### Step 9: Install and Test the Formula

```bash
# Install your formula from the tap
brew install vib795/tap/pull-vids

# Test that it works
pull-vids --version
pull-vids --help

# Run Homebrew's test block
brew test vib795/tap/pull-vids
```

**Expected output:**
```
✅ pull-vids --version returns version info
✅ pull-vids --help shows "Universal Video Downloader"
✅ All tests pass
```

---

### Step 10: Run Homebrew Audit

```bash
# Run audit to check for issues
brew audit --strict vib795/tap/pull-vids

# For new formula specifically
brew audit --new vib795/tap/pull-vids
```

**Expected output:**
```
✅ No errors
✅ No warnings
✅ Formula follows Homebrew style guide
```

**Common warnings and fixes:**
- **Checksum mismatch**: Re-download from GitHub and verify checksums match
- **Missing dependencies**: Ensure `ffmpeg` and `yt-dlp` are in homebrew-core
- **URL issues**: Verify release v0.2.0 exists and is public

---

### Step 11: Uninstall and Reinstall Test

```bash
# Uninstall
brew uninstall pull-vids

# Reinstall to verify it works cleanly
brew install vib795/tap/pull-vids

# Test again
pull-vids --version
```

---

## 📋 Homebrew Core Contribution Checklist

Now you can answer the checklist with confidence:

### ✅ 1. Have you built your formula locally?
**Answer:** Yes
```bash
brew install vib795/tap/pull-vids
```

### ✅ 2. Is your test running fine?
**Answer:** Yes
```bash
brew test vib795/tap/pull-vids
```

### ✅ 3. Does your build pass audit?
**Answer:** Yes
```bash
brew audit --strict vib795/tap/pull-vids
brew audit --new vib795/tap/pull-vids
```

### ✅ 4. Does your formula follow Homebrew style guide?
**Answer:** Yes
- Uses `on_macos` and `on_linux` blocks
- Uses `on_arm` and `on_intel` for architecture
- Includes `desc`, `homepage`, `license`
- Has proper dependencies with `depends_on`
- Test block validates version and help output

### ✅ 5. Is this your first contribution to Homebrew?
**Answer:** [Update based on your situation]
- If yes, read: https://docs.brew.sh/How-To-Open-a-Homebrew-Pull-Request

### ✅ 6. Have you ensured there are no other open PRs for this formula?
**Answer:** Yes
- Check: https://github.com/Homebrew/homebrew-core/pulls?q=is%3Apr+pull-vids

---

## 🎯 Submit to Homebrew Core

Once all tests pass, submit to Homebrew core:

### Step 1: Fork homebrew-core

```bash
# Fork the repository at:
# https://github.com/Homebrew/homebrew-core

# Clone your fork
cd ~/
git clone https://github.com/YOUR_USERNAME/homebrew-core.git
cd homebrew-core
```

### Step 2: Create a Branch

```bash
# Create a branch for your formula
git checkout -b pull-vids
```

### Step 3: Add Your Formula

```bash
# Copy your formula to the correct location
# Note: Formulas are organized alphabetically in subdirectories
cp ~/i/pull-vids/Formula/pull-vids.rb Formula/p/pull-vids.rb

# Verify it's there
cat Formula/p/pull-vids.rb
```

### Step 4: Commit with Proper Message

```bash
# Stage the formula
git add Formula/p/pull-vids.rb

# Commit with Homebrew's required format
git commit -m "pull-vids 0.2.0 (new formula)"
```

**Important:** The commit message must follow this exact format:
- For new formulas: `pull-vids 0.2.0 (new formula)`
- For updates: `pull-vids 0.2.1`

### Step 5: Push to Your Fork

```bash
git push origin pull-vids
```

### Step 6: Create Pull Request

1. Go to: https://github.com/Homebrew/homebrew-core/compare
2. Click "compare across forks"
3. Select:
   - **Base repository:** `Homebrew/homebrew-core`
   - **Base branch:** `master`
   - **Head repository:** `YOUR_USERNAME/homebrew-core`
   - **Compare branch:** `pull-vids`
4. Title: `pull-vids 0.2.0 (new formula)`
5. Fill out the checklist in the PR template
6. Submit!

---

## 🧪 Testing Commands Reference

```bash
# Install formula
brew install vib795/tap/pull-vids

# Uninstall formula
brew uninstall pull-vids

# Test formula
brew test vib795/tap/pull-vids

# Audit formula
brew audit --strict vib795/tap/pull-vids
brew audit --new vib795/tap/pull-vids

# View formula info
brew info vib795/tap/pull-vids

# View formula source
brew cat vib795/tap/pull-vids

# List installed versions
brew list --versions pull-vids
```

---

## 🔧 Troubleshooting

### Issue: "No available formula with the name pull-vids"
**Fix:** You must specify the tap or file path:
```bash
# Wrong
brew install pull-vids

# Right
brew install vib795/tap/pull-vids
# Or
brew install Formula/pull-vids.rb
```

### Issue: "SHA256 mismatch"
**Fix:** Checksums don't match the downloaded file
```bash
# 1. Delete the cached download
rm -rf ~/Library/Caches/Homebrew/downloads/*pull-vids*

# 2. Re-calculate checksums
cd dist
shasum -a 256 *.tar.gz

# 3. Update Formula/pull-vids.rb with new checksums
# 4. Update the tap
cp Formula/pull-vids.rb $(brew --repository)/Library/Taps/vib795/homebrew-tap/Formula/

# 5. Try again
brew reinstall vib795/tap/pull-vids
```

### Issue: "Release already exists: v0.2.0"
**Fix:** Delete and recreate the release
```bash
gh release delete v0.2.0 --yes
# Then re-run Step 7
```

### Issue: "Dependencies not satisfied"
**Fix:** Install required dependencies
```bash
brew install ffmpeg
brew install yt-dlp
```

### Issue: "Formula requires formulae to be in a tap"
**Fix:** Don't install directly from file, use a tap:
```bash
# Wrong
brew install Formula/pull-vids.rb

# Right - create tap first
brew tap-new vib795/tap
cp Formula/pull-vids.rb $(brew --repository)/Library/Taps/vib795/homebrew-tap/Formula/
brew install vib795/tap/pull-vids
```

---

## 📝 Quick Checklist

Before submitting to Homebrew core, verify:

- [ ] Binaries built for all platforms
- [ ] SHA256 checksums calculated and verified
- [ ] Formula updated with correct checksums
- [ ] GitHub release v0.2.0 created with all artifacts
- [ ] Formula installs successfully: `brew install vib795/tap/pull-vids`
- [ ] Tests pass: `brew test vib795/tap/pull-vids`
- [ ] Audit passes: `brew audit --new vib795/tap/pull-vids`
- [ ] Binary works: `pull-vids --version` and `pull-vids --help`
- [ ] No other open PRs for pull-vids
- [ ] Commit message follows format: `pull-vids 0.2.0 (new formula)`

---

## 📚 Resources

- **Homebrew Formula Cookbook:** https://docs.brew.sh/Formula-Cookbook
- **Homebrew PR Guide:** https://docs.brew.sh/How-To-Open-a-Homebrew-Pull-Request
- **Homebrew Style Guide:** https://docs.brew.sh/Formula-Cookbook#style-guide
- **Your GitHub Release:** https://github.com/vib795/pull-vids/releases/tag/v0.2.0
- **Homebrew Core:** https://github.com/Homebrew/homebrew-core

---

## 🎉 Success!

Once your PR is merged:
- Your formula will be available as: `brew install pull-vids`
- Users worldwide can install your tool with a single command
- You're now a Homebrew contributor!

Good luck! 🍺
