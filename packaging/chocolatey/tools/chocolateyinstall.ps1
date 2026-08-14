$ErrorActionPreference = 'Stop'

# $version$ and $checksum$ are substituted at pack time by the release workflow,
# so the checksum always matches the asset built for this exact version.
$packageName    = 'pull-vids'
$packageVersion = '$version$'
$url64          = "https://github.com/vib795/pull-vids/releases/download/v$packageVersion/pull-vids-windows-amd64.zip"
$checksum64     = '$checksum$'

$toolsDir = "$(Split-Path -Parent $MyInvocation.MyCommand.Definition)"

$packageArgs = @{
  packageName    = $packageName
  unzipLocation  = $toolsDir
  url64bit       = $url64
  checksum64     = $checksum64
  checksumType64 = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs

# The archive holds pull-vids-windows-amd64.exe. Chocolatey shims every .exe it
# finds, which would expose the command under that name, so rename it to the
# documented command name before shimming.
$original = Join-Path $toolsDir 'pull-vids-windows-amd64.exe'
$target   = Join-Path $toolsDir 'pull-vids.exe'

if (Test-Path $original) {
  Move-Item -Path $original -Destination $target -Force
}

if (-not (Test-Path $target)) {
  throw "pull-vids.exe was not found in $toolsDir after extraction."
}
