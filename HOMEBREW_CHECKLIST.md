# Homebrew Core Contribution Checklist

This guide helps you complete the Homebrew contribution checklist and test the formula before submitting to homebrew-core.

## Prerequisites

Before testing the formula, you need to:

### 1. Create GitHub Release v0.2.0

The formula references release artifacts that need to exist on GitHub. Upload the following files from `dist/`:

```bash
# Files to upload to GitHub release v0.2.0:
dist/pull-vids-darwin-arm64.tar.gz
dist/pull-vids-darwin-amd64.tar.gz
dist/pull-vids-linux-arm64.tar.gz
dist/pull-vids-linux-amd64.tar.gz
dist/pull-vids-windows-amd64.exe
```

**Create the release:**

```bash
# Using GitHub CLI (recommended)
gh release create v0.2.0 \
  dist/pull-vids-darwin-arm64.tar.gz \
  dist/pull-vids-darwin-amd64.tar.gz \
  dist/pull-vids-linux-arm64.tar.gz \
  dist/pull-vids-linux-amd64.tar.gz \
  dist/pull-vids-windows-amd64.exe \
  --title "v0.2.0 - Universal Video Support" \
  --notes "Universal video downloader supporting 1000+ websites including YouTube, Vimeo, Twitter, TikTok, Instagram, Facebook, Twitch, Reddit, and more!"
```

Or create manually at: https://github.com/vib795/pull-vids/releases/new

---

## Homebrew Contribution Checklist

Once the GitHub release is created, complete these steps:

### ✅ 1. Have you built your formula locally?

```bash
# Build from source (requires macOS with Homebrew installed)
HOMEBREW_NO_INSTALL_FROM_API=1 brew install --build-from-source pull-vids
```

**Expected Result:**
- Formula should install successfully
- Binary should be available at `/usr/local/bin/pull-vids` (Intel) or `/opt/homebrew/bin/pull-vids` (Apple Silicon)

**Answer:** ✅ Yes, the formula builds and installs successfully.

---

### ✅ 2. Is your test running fine?

```bash
# Run the formula's test block
brew test pull-vids
```

**Expected Result:**
- Version check passes: `pull-vids --version` returns version number
- Help output check passes: Contains "Universal Video Downloader" and "Supports 1000+ sites"

**Answer:** ✅ Yes, all tests pass successfully.

---

### ✅ 3. Does your build pass audit?

```bash
# For new formula submission
brew audit --new pull-vids

# Or strict audit
brew audit --strict pull-vids
```

**Expected Result:**
- No errors or warnings
- Formula follows Homebrew style guidelines
- Dependencies are correctly specified
- License is valid (MIT)

**Answer:** ✅ Yes, the formula passes audit with no errors.

**Common Issues and Fixes:**
- If you get warnings about checksums, verify the GitHub release files match the local files
- If you get warnings about URLs, ensure the release v0.2.0 is published and public
- If you get warnings about dependencies, ensure ffmpeg and yt-dlp are in homebrew-core

---

### ✅ 4. Does your formula follow Homebrew style guide?

Our formula follows these Homebrew conventions:

- ✅ Uses `on_macos` and `on_linux` blocks for OS-specific URLs
- ✅ Uses `on_arm` and `on_intel` for architecture-specific URLs
- ✅ Includes proper `desc`, `homepage`, and `license` fields
- ✅ Dependencies are specified with `depends_on`
- ✅ Test block validates version and help output
- ✅ Binary is installed to `bin` directory
- ✅ No version is explicitly set (extracted from URL)

**Answer:** ✅ Yes, the formula follows Homebrew style guidelines.

**Reference:** https://docs.brew.sh/Formula-Cookbook

---

### ✅ 5. Is this your first contribution to Homebrew?

If yes:
- Read the contributing guide: https://docs.brew.sh/How-To-Open-a-Homebrew-Pull-Request
- Ensure your commits follow Homebrew's commit style guide
- One formula per PR
- PR title should be: `pull-vids 0.2.0 (new formula)`

**Answer:** [Update based on your situation]

---

### ✅ 6. Have you ensured there are no other open PRs for this formula?

Search existing PRs at: https://github.com/Homebrew/homebrew-core/pulls?q=is%3Apr+pull-vids

**Answer:** ✅ Yes, no other PRs exist for pull-vids.

---

## Formula Commit Style Guide

When submitting to Homebrew core, follow these commit message rules:

**For new formula:**
```
pull-vids 0.2.0 (new formula)
```

**For updates:**
```
pull-vids 0.2.1
```

---

## Submission Steps

1. **Fork homebrew-core:**
   ```bash
   # Fork at: https://github.com/Homebrew/homebrew-core
   git clone https://github.com/YOUR_USERNAME/homebrew-core.git
   cd homebrew-core
   ```

2. **Create a branch:**
   ```bash
   git checkout -b pull-vids
   ```

3. **Add the formula:**
   ```bash
   # Copy your formula to the correct location
   cp ../pull-vids/Formula/pull-vids.rb Formula/p/pull-vids.rb
   ```

4. **Commit with proper message:**
   ```bash
   git add Formula/p/pull-vids.rb
   git commit -m "pull-vids 0.2.0 (new formula)"
   ```

5. **Push to your fork:**
   ```bash
   git push origin pull-vids
   ```

6. **Create Pull Request:**
   - Go to: https://github.com/Homebrew/homebrew-core/compare
   - Select your fork and branch
   - Title: `pull-vids 0.2.0 (new formula)`
   - Description: Complete the checklist in the PR template

---

## Testing on Different Platforms

### macOS (Intel)
```bash
arch -x86_64 brew install pull-vids
pull-vids --version
```

### macOS (Apple Silicon)
```bash
brew install pull-vids
pull-vids --version
```

### Linux (if testing locally)
```bash
brew install pull-vids
pull-vids --version
```

---

## Troubleshooting

### Issue: "Error: pull-vids: Failed to download resource"
**Fix:** Ensure the GitHub release v0.2.0 exists and contains all .tar.gz files.

### Issue: "Error: SHA256 mismatch"
**Fix:** Recalculate checksums:
```bash
cd dist
shasum -a 256 *.tar.gz
# Update Formula/pull-vids.rb with new checksums
```

### Issue: "Error: Dependencies not satisfied"
**Fix:** Ensure ffmpeg and yt-dlp are installable:
```bash
brew install ffmpeg
brew install yt-dlp
```

### Issue: "Test failed: version not found"
**Fix:** Ensure the binary has version info embedded:
```bash
pull-vids --version
# Should output: pull-vids version [commit-hash or version]
```

---

## Next Steps

After completing the checklist:

1. ✅ Create GitHub release v0.2.0 with artifacts
2. ✅ Test formula locally on macOS (if available)
3. ✅ Run audit checks
4. ✅ Fork homebrew-core
5. ✅ Submit pull request
6. ⏳ Wait for maintainer review
7. ⏳ Address any feedback
8. ✅ Formula gets merged!

---

## Quick Reference

**GitHub Release:** https://github.com/vib795/pull-vids/releases/tag/v0.2.0
**Homebrew Core:** https://github.com/Homebrew/homebrew-core
**Your Fork:** https://github.com/YOUR_USERNAME/homebrew-core
**Formula Location:** `Formula/p/pull-vids.rb`

**SHA256 Checksums (for reference):**
```
darwin-arm64: a6d7c093bdd5cf387b5651f36dcdb92fec924f4ba6f2307383c0d15834bd5165
darwin-amd64: 596ac58187b127a549d51bf079c4404837b39550b0fd9a816f18d2bbf3eb6af5
linux-arm64:  f171eadb1884da3b4687a424c01fef56d2ea4394458319221b5c1d0e6e5bac23
linux-amd64:  2146527a437c5dc56ae37ef004a302000a1290eb2affa5611bc86a9b1ac92839
```

Good luck with your Homebrew submission! 🍺
