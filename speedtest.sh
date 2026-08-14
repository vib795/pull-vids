#!/usr/bin/env bash
# speedtest.sh — compare pull-vids download throughput across connection settings.
#
#   ./speedtest.sh                 # use the built-in sample video
#   ./speedtest.sh "<youtube-url>" # use your own
#
# Each configuration downloads to a scratch directory, is timed end to end, and
# reports effective MiB/s computed from the bytes actually written. The scratch
# directory is removed afterwards.

set -uo pipefail
cd "$(dirname "$0")"

URL="${1:-https://www.youtube.com/watch?v=aqz-KE-bpKQ}"
QUALITY="${QUALITY:-1080p}"
OLD="./pull-vids"
NEW="./pull-vids-fast"

if [ ! -x "$NEW" ]; then
  echo "Building $NEW ..."
  go build -o "$NEW" . || exit 1
fi

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

human() { awk -v b="$1" 'BEGIN{printf "%.1f MiB", b/1048576}'; }

# run <label> <binary> [extra args...]
run() {
  local label="$1" bin="$2"; shift 2
  local dir="$WORK/$RANDOM$RANDOM"
  mkdir -p "$dir"

  local t0 t1 bytes secs rate rc
  t0=$(python3 -c 'import time;print(time.time())')
  "$bin" --no-banner -q "$QUALITY" -o "$dir" "$@" "$URL" >"$dir/.log" 2>&1
  rc=$?
  t1=$(python3 -c 'import time;print(time.time())')

  bytes=$(find "$dir" -type f ! -name '.log' -exec stat -f%z {} + 2>/dev/null | awk '{s+=$1} END{print s+0}')

  if [ "$rc" -ne 0 ] || [ "$bytes" -eq 0 ]; then
    printf '  %-22s FAILED  %s\n' "$label" \
      "$(grep -oiE 'ERROR:.*' "$dir/.log" | head -1 | cut -c1-70)"
    rm -rf "$dir"; return
  fi

  secs=$(python3 -c "print(f'{$t1-$t0:.2f}')")
  rate=$(python3 -c "print(f'{$bytes/1048576/($t1-$t0):.2f}')")
  printf '  %-22s %8ss  %10s  %8s MiB/s\n' "$label" "$secs" "$(human "$bytes")" "$rate"
  rm -rf "$dir"
}

echo
echo "URL:     $URL"
echo "Quality: $QUALITY"
if command -v aria2c >/dev/null 2>&1; then
  echo "aria2c:  installed ($(aria2c --version | head -1 | awk '{print $3}'))"
else
  echo "aria2c:  NOT installed  ->  brew install aria2   (unlocks the biggest gain)"
fi
echo
printf '  %-22s %9s  %10s  %13s\n' "CONFIG" "TIME" "SIZE" "THROUGHPUT"
printf '  %s\n' "------------------------------------------------------------------"

# Old binary: single stream, no tuning flags. This is the current behaviour.
[ -x "$OLD" ] && run "current (1 conn)" "$OLD"

run "new: native -N 4"  "$NEW" --downloader native -N 4
run "new: native -N 8"  "$NEW" --downloader native -N 8
run "new: native -N 16" "$NEW" --downloader native -N 16

if command -v aria2c >/dev/null 2>&1; then
  run "new: aria2c -N 8"  "$NEW" --downloader aria2c -N 8
  run "new: aria2c -N 16" "$NEW" --downloader aria2c -N 16
fi

echo
echo "Note: YouTube throttles per connection and per IP. If you see repeated"
echo "      failures, wait a few minutes before re-running - back-to-back runs"
echo "      can trip rate limiting and skew the comparison."
echo
