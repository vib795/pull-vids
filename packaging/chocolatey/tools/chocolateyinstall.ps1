$ErrorActionPreference = 'Stop'

$packageName = 'pull-vids'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-windows-amd64.zip'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  url64bit      = $url64
  checksum64    = 'REPLACE_WITH_ACTUAL_CHECKSUM'
  checksumType64= 'sha256'
}

Install-ChocolateyZipPackage @packageArgs

# Rename binary
$exePath = Join-Path $toolsDir 'pull-vids-windows-amd64.exe'
$targetPath = Join-Path $toolsDir 'pull-vids.exe'
if (Test-Path $exePath) {
  Rename-Item -Path $exePath -NewName 'pull-vids.exe' -Force
}

Write-Host ""
Write-Host "pull-vids has been installed!" -ForegroundColor Green
Write-Host ""
Write-Host "Note: You also need to install:" -ForegroundColor Yellow
Write-Host "  - yt-dlp: pip install yt-dlp" -ForegroundColor Cyan
Write-Host ""
Write-Host "Run 'pull-vids --help' to get started!" -ForegroundColor Green
