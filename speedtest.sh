#!/usr/bin/env bash
# speedtest.sh — compare pull-vids download throughput between its backends.
#
#   ./speedtest.sh                 # use the built-in sample video
#   ./speedtest.sh "<youtube-url>" # use your own
#   QUALITY=720p ./speedtest.sh    # pick another quality (default 1080p)
#
# The current source is built fresh, so the numbers always describe the code in
# this checkout. Each configuration downloads to a scratch directory, is timed
# end to end, and reports effective MiB/s from the bytes actually written.

set -uo pipefail
cd "$(dirname "$0")"

URL="${1:-https://www.youtube.com/watch?v=dQw4w9WgXcQ}"
QUALITY="${QUALITY:-1080p}"

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

BIN="$WORK/pull-vids"
echo "Building $BIN ..."
go build -o "$BIN" . || exit 1

now() { python3 -c 'import time; print(time.time())'; }

# Sums the sizes of the downloaded files. stat's flags differ between macOS and
# Linux, so python3, which the timing already needs, does it portably.
dir_bytes() {
  python3 - "$1" <<'EOF'
import os, sys
print(sum(os.path.getsize(os.path.join(d, f))
          for d, _, files in os.walk(sys.argv[1]) for f in files if f != '.log'))
EOF
}

# run <label> [extra args...]
run() {
  local label="$1"; shift
  local dir="$WORK/$RANDOM$RANDOM"
  mkdir -p "$dir"

  local t0 t1 bytes rc
  t0=$(now)
  "$BIN" --no-banner -q "$QUALITY" -o "$dir" "$@" "$URL" >"$dir/.log" 2>&1
  rc=$?
  t1=$(now)
  bytes=$(dir_bytes "$dir")

  if [ "$rc" -ne 0 ] || [ "$bytes" -eq 0 ]; then
    printf '  %-22s FAILED  %s\n' "$label" \
      "$(grep -oiE 'ERROR:.*' "$dir/.log" | head -1 | cut -c1-70)"
  else
    python3 -c "
b, s = $bytes, $t1 - $t0
print(f'  {\"$label\":<22} {s:8.2f}s  {b/1048576:6.1f} MiB  {b/1048576/s:8.2f} MiB/s')"
  fi
  rm -rf "$dir"
}

echo
echo "URL:     $URL"
echo "Quality: $QUALITY"
echo
printf '  %-22s %9s  %10s  %14s\n' "CONFIG" "TIME" "SIZE" "THROUGHPUT"
printf '  %s\n' "------------------------------------------------------------"

run "native (default)"
if command -v aria2c >/dev/null 2>&1; then
  run "aria2c -N 8" --downloader aria2c -N 8
else
  echo "  aria2c                 skipped (not installed)"
fi

echo
echo "Note: back-to-back runs can trip YouTube's rate limiting and skew the"
echo "      comparison. If a run fails, wait a few minutes and try again."
echo
