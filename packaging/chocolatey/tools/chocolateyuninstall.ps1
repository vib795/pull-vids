$ErrorActionPreference = 'Stop'

$packageName = 'pull-vids'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

# Remove binary
$exePath = Join-Path $toolsDir 'pull-vids.exe'
if (Test-Path $exePath) {
  Remove-Item -Path $exePath -Force
}

Write-Host "pull-vids has been uninstalled!" -ForegroundColor Green
