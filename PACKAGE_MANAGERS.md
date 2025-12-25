# Package Manager Distribution Guide

This guide explains how to publish pull-vids to various package managers.

## Homebrew (macOS/Linux)

### Option 1: Homebrew Tap (Recommended)

Create your own tap repository:

```bash
# Create a tap repository
brew tap-new vib795/pull-vids

# Copy the formula
cp Formula/pull-vids.rb $(brew --repository)/Library/Taps/vib795/homebrew-pull-vids/Formula/

# Test the formula
brew install vib795/pull-vids/pull-vids
```

**Users can install with:**
```bash
brew tap vib795/pull-vids
brew install pull-vids
```

### Option 2: Homebrew Core (Official)

To add to Homebrew core:

1. Fork https://github.com/Homebrew/homebrew-core
2. Add `Formula/pull-vids.rb` to the repository
3. Submit a pull request
4. Wait for review and merge

**Users can install with:**
```bash
brew install pull-vids
```

### Updating the Formula

After each release:

```bash
# Calculate SHA256 checksums
shasum -a 256 dist/releases/*.tar.gz

# Update Formula/pull-vids.rb with:
# - New version number
# - New URL
# - New SHA256 checksums
```

## Chocolatey (Windows)

### Setup

1. Create a Chocolatey account at https://community.chocolatey.org/
2. Get your API key from your account settings

### Publishing

```powershell
# Navigate to packaging directory
cd packaging/chocolatey

# Update version and checksums in pull-vids.nuspec
# Calculate checksum:
# (Get-FileHash -Path pull-vids-windows-amd64.zip -Algorithm SHA256).Hash

# Pack the package
choco pack

# Test locally
choco install pull-vids -s . -y

# Push to Chocolatey
choco push pull-vids.0.2.0.nupkg --source https://push.chocolatey.org/ --api-key YOUR_API_KEY
```

**Users can install with:**
```powershell
choco install pull-vids
```

## APT (Debian/Ubuntu)

### Creating .deb Package

```bash
# Build for Linux AMD64
make build-all

# Create package structure
mkdir -p pull-vids_0.2.0_amd64/usr/local/bin
mkdir -p pull-vids_0.2.0_amd64/DEBIAN

# Copy files
cp dist/pull-vids-linux-amd64 pull-vids_0.2.0_amd64/usr/local/bin/pull-vids
cp packaging/deb/DEBIAN/* pull-vids_0.2.0_amd64/DEBIAN/
chmod +x pull-vids_0.2.0_amd64/DEBIAN/postinst
chmod 755 pull-vids_0.2.0_amd64/usr/local/bin/pull-vids

# Build the package
dpkg-deb --build pull-vids_0.2.0_amd64

# Test installation
sudo dpkg -i pull-vids_0.2.0_amd64.deb
```

### Publishing to PPA (Ubuntu)

1. Create a Launchpad account
2. Set up GPG keys
3. Create PPA: https://launchpad.net/~USERNAME/+activate-ppa
4. Upload package:

```bash
dput ppa:USERNAME/pull-vids pull-vids_0.2.0_amd64.changes
```

**Users can install with:**
```bash
sudo add-apt-repository ppa:vib795/pull-vids
sudo apt update
sudo apt install pull-vids
```

## Scoop (Windows - Alternative)

Create a bucket repository:

```json
{
    "version": "0.2.0",
    "description": "Universal video downloader CLI supporting 1000+ websites",
    "homepage": "https://github.com/vib795/pull-vids",
    "license": "MIT",
    "url": "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-windows-amd64.zip",
    "hash": "REPLACE_WITH_SHA256",
    "bin": "pull-vids.exe",
    "depends": ["ffmpeg"],
    "suggest": {
        "yt-dlp": "python install yt-dlp"
    },
    "checkver": "github",
    "autoupdate": {
        "url": "https://github.com/vib795/pull-vids/releases/download/v$version/pull-vids-windows-amd64.zip"
    }
}
```

**Users can install with:**
```powershell
scoop bucket add pull-vids https://github.com/vib795/scoop-pull-vids
scoop install pull-vids
```

## Snap (Linux - Universal)

Create `snap/snapcraft.yaml`:

```yaml
name: pull-vids
version: '0.2.0'
summary: Universal video downloader CLI supporting 1000+ websites
description: |
  A free, open-source CLI tool for downloading videos and audio from 1000+ websites.
  Inspired by tools like Downie and PullTube, but completely free!
  .
  Supports YouTube, Vimeo, Twitter, TikTok, Instagram, Facebook, Twitch, Reddit, and more!

base: core20
confinement: strict
grade: stable

apps:
  pull-vids:
    command: bin/pull-vids
    plugs: [home, network]

parts:
  pull-vids:
    plugin: go
    source: .
    build-packages:
      - gcc
```

**Build and publish:**
```bash
snapcraft
snapcraft login
snapcraft upload --release=stable pull-vids_0.2.0_amd64.snap
```

**Users can install with:**
```bash
sudo snap install pull-vids
```

## AUR (Arch Linux)

Create `PKGBUILD`:

```bash
pkgname=pull-vids
pkgver=0.2.0
pkgrel=1
pkgdesc="Universal video downloader CLI supporting 1000+ websites"
arch=('x86_64' 'aarch64')
url="https://github.com/vib795/pull-vids"
license=('MIT')
depends=('ffmpeg' 'yt-dlp')
source_x86_64=("$pkgname-$pkgver-x86_64.tar.gz::https://github.com/vib795/pull-vids/releases/download/v$pkgver/pull-vids-linux-amd64.tar.gz")
source_aarch64=("$pkgname-$pkgver-aarch64.tar.gz::https://github.com/vib795/pull-vids/releases/download/v$pkgver/pull-vids-linux-arm64.tar.gz")
sha256sums_x86_64=('REPLACE_WITH_SHA256')
sha256sums_aarch64=('REPLACE_WITH_SHA256')

package() {
    install -Dm755 "$srcdir/pull-vids-linux-amd64" "$pkgdir/usr/bin/pull-vids"
}
```

**Users can install with:**
```bash
yay -S pull-vids
# or
paru -S pull-vids
```

## Nix (NixOS)

Create a derivation and submit to nixpkgs:

```nix
{ lib, buildGoModule, fetchFromGitHub, ffmpeg, yt-dlp }:

buildGoModule rec {
  pname = "pull-vids";
  version = "0.2.0";

  src = fetchFromGitHub {
    owner = "vib795";
    repo = "pull-vids";
    rev = "v${version}";
    sha256 = "REPLACE_WITH_SHA256";
  };

  vendorSha256 = "REPLACE_WITH_VENDOR_SHA256";

  propagatedBuildInputs = [ ffmpeg yt-dlp ];

  meta = with lib; {
    description = "Universal video downloader CLI supporting 1000+ websites";
    homepage = "https://github.com/vib795/pull-vids";
    license = licenses.mit;
    maintainers = with maintainers; [ ];
  };
}
```

**Users can install with:**
```bash
nix-env -iA nixpkgs.pull-vids
```

## Automation

Add to `Makefile`:

```makefile
publish-brew: release
	@echo "Publishing to Homebrew..."
	# Update checksums in Formula/pull-vids.rb
	# Push to tap repository

publish-choco: release
	@echo "Publishing to Chocolatey..."
	cd packaging/chocolatey && choco pack && choco push

publish-deb: release
	@echo "Creating .deb package..."
	# Build .deb package
	# Upload to PPA

publish-all: publish-brew publish-choco publish-deb
	@echo "Published to all package managers!"
```

## Summary

Once set up, users can install with:

```bash
# macOS
brew install pull-vids

# Windows
choco install pull-vids
# or
scoop install pull-vids

# Ubuntu/Debian
sudo apt install pull-vids

# Arch Linux
yay -S pull-vids

# NixOS
nix-env -iA nixpkgs.pull-vids

# Universal (Linux)
sudo snap install pull-vids
```
