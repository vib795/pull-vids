# Install script for pull-vids on Windows
# Run this in PowerShell as Administrator

$ErrorActionPreference = "Stop"

$REPO = "vib795/pull-vids"
$BINARY_NAME = "pull-vids.exe"
$INSTALL_DIR = "$env:ProgramFiles\pull-vids"

Write-Host "╔═══════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║     pull-vids Installation Script     ║" -ForegroundColor Cyan
Write-Host "╚═══════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Check if running as Administrator
$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
if (-not $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "✗ This script must be run as Administrator" -ForegroundColor Red
    Write-Host "  Right-click PowerShell and select 'Run as Administrator'" -ForegroundColor Yellow
    exit 1
}

Write-Host "Detected platform: Windows AMD64" -ForegroundColor Yellow
Write-Host ""

# Get latest release
Write-Host "Fetching latest release..." -ForegroundColor Cyan
try {
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
    $version = $release.tag_name
    Write-Host "✓ Latest release: $version" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed to fetch latest release" -ForegroundColor Red
    Write-Host "  Building from source not yet supported on Windows" -ForegroundColor Yellow
    exit 1
}

# Download and unpack the release archive.
# Releases ship a zip, not a bare .exe, so fetching "$BINARY_NAME" directly
# returns 404 and the install fails.
$ASSET_NAME = "pull-vids-windows-amd64.zip"
$EXE_IN_ZIP = "pull-vids-windows-amd64.exe"

$downloadUrl = "https://github.com/$REPO/releases/download/$version/$ASSET_NAME"
$tempDir = [System.IO.Path]::Combine([System.IO.Path]::GetTempPath(), "pull-vids-$version")
$tempZip = [System.IO.Path]::Combine([System.IO.Path]::GetTempPath(), $ASSET_NAME)

Write-Host "Downloading $ASSET_NAME..." -ForegroundColor Cyan
try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempZip
    Write-Host "✓ Download complete" -ForegroundColor Green
} catch {
    Write-Host "✗ Download failed: $downloadUrl" -ForegroundColor Red
    exit 1
}

Write-Host "Extracting..." -ForegroundColor Cyan
try {
    if (Test-Path $tempDir) { Remove-Item -Path $tempDir -Recurse -Force }
    Expand-Archive -Path $tempZip -DestinationPath $tempDir -Force

    $tempFile = [System.IO.Path]::Combine($tempDir, $EXE_IN_ZIP)
    if (-not (Test-Path $tempFile)) {
        # Fall back to whatever .exe the archive contains, in case the asset
        # naming changes in a future release.
        $found = Get-ChildItem -Path $tempDir -Filter "*.exe" -Recurse | Select-Object -First 1
        if (-not $found) { throw "no .exe found inside $ASSET_NAME" }
        $tempFile = $found.FullName
    }
    Write-Host "✓ Extracted" -ForegroundColor Green
} catch {
    Write-Host "✗ Extraction failed: $_" -ForegroundColor Red
    exit 1
}

# Create installation directory
Write-Host "Installing to $INSTALL_DIR..." -ForegroundColor Cyan
if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR | Out-Null
}

# Copy binary
Copy-Item -Path $tempFile -Destination "$INSTALL_DIR\$BINARY_NAME" -Force
Remove-Item -Path $tempZip -Force -ErrorAction SilentlyContinue
Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue

# Add to PATH if not already there
$currentPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::Machine)
if ($currentPath -notlike "*$INSTALL_DIR*") {
    Write-Host "Adding to system PATH..." -ForegroundColor Cyan
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$currentPath;$INSTALL_DIR",
        [EnvironmentVariableTarget]::Machine
    )
    Write-Host "✓ Added to PATH" -ForegroundColor Green
}

Write-Host ""
Write-Host "╔═══════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║  ✓ Installation completed!            ║" -ForegroundColor Green
Write-Host "╚═══════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""
Write-Host "Run 'pull-vids --help' to get started!" -ForegroundColor Yellow
Write-Host "(You may need to restart your terminal)" -ForegroundColor Yellow
Write-Host ""
Write-Host "Note: You also need to install:" -ForegroundColor Cyan
Write-Host "  - ffmpeg (for video processing)" -ForegroundColor White
Write-Host "  - yt-dlp (for downloading from 1000+ sites)" -ForegroundColor White
Write-Host ""
Write-Host "Install dependencies:" -ForegroundColor Cyan
Write-Host "  winget install ffmpeg" -ForegroundColor Yellow
Write-Host "  pip install yt-dlp" -ForegroundColor Yellow
Write-Host ""
