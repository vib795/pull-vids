# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A single-file Go CLI (`main.go`, ~576 lines, `package main`) that wraps the `yt-dlp` binary. There are no internal packages and no abstractions — the design is intentionally flat. The program's real job is threefold: translate friendly flags into `yt-dlp` arguments, parse the subprocess's stdout into a live progress bar, and retry intelligently when a CDN pushes back.

Runtime dependencies are external binaries, not Go libraries: `yt-dlp` (checked at startup by `checkYtDlp()`), `ffmpeg` (required by yt-dlp, never checked), and `aria2c` (optional, auto-detected).

## Commands

```bash
make build        # → ./pull-vids, with version injected from git describe
make build-all    # cross-compile 5 targets → dist/
make release      # wrap dist/ binaries into .tar.gz / .zip → dist/releases/
make checksums    # release + SHA256SUMS
make install      # sudo cp to /usr/local/bin
make help         # lists all targets (this is .DEFAULT_GOAL)
```

**There are no tests.** `make test` runs `go test -v ./...` against zero `_test.go` files and trivially passes — do not read a green `make test` as verification of anything. Changes are verified by running the binary:

```bash
./pull-vids -q 720p "https://www.youtube.com/watch?v=..."
./pull-vids -N 16 --downloader aria2c "<url>"   # exercise the parallel path
./pull-vids --downloader native -N 8 "<url>"    # exercise the fallback path
./speedtest.sh                                   # benchmark both backends across -N values
```

`make deb` is broken — it reads `packaging/deb/DEBIAN/*`, which does not exist in the repo.

## Throughput architecture

This is the part most likely to be misunderstood. Google's CDN rate-limits **each TCP connection independently** at roughly 3 MB/s, so a single-stream download leaves a gigabit link idle. Aggregate throughput scales close to linearly with connection count (measured: 1→3.3, 8→25.8, 16→48.7, 32→95.2 MB/s). Every speed flag exists to buy more connections.

`-N/-connections` (default 8) means different things depending on the backend selected by `resolveDownloader()`:

- **aria2c** (preferred when installed) → `--downloader-args aria2c:-xN -sN -jN -k1M ...`
  - `-x`/`-s` split one contiguous URL into parallel ranged requests.
  - **`-j` is separate and equally critical.** For DASH/HLS formats yt-dlp hands aria2c a *fragment list*; without `-j` aria2c fetches those fragments serially, which is **slower than the native downloader**. This was shipped broken in v0.3.0 and fixed in v0.3.4 — do not drop `-j` when editing these args.
  - `-k1M` is required because aria2c refuses to split ranges under 20M by default.
- **native** → `--concurrent-fragments N` plus `--http-chunk-size`. This only parallelises formats that are genuinely fragmented.

`downloadVideo()` wraps `executeDownload()` in a retry loop that treats **HTTP 403 as a concurrency signal, not a dead URL** — it halves `config.Connections` before each retry rather than retrying identically.

Because the backend is switchable, **two progress formats must be parsed** in the stdout goroutine: yt-dlp's `[download] 45.3% of 12.34MiB at 1.23MiB/s ETA 00:10` lines, and aria2c's `[#8a1b2c 12MiB/100MiB(12%) CN:16 DL:25MiB ETA:3s]` (via the `aria2Progress` regex). Adding a backend means adding a third parser.

## Traps

**`version` must stay a `var`.** The linker cannot patch a `const`, so declaring it const makes `-ldflags "-X main.version=..."` a silent no-op that ships the hardcoded fallback. This shipped as a real bug once.

**The Homebrew formula needs an explicit `version` stanza.** Without it Homebrew infers the version from the URL filename and reads `arm64` as version `"64"`. Since `Version.new("64") > Version.new("0.3.5")`, affected installs then refuse to upgrade forever and need `brew uninstall` + `brew install`. The stanza is generated in `update-homebrew.yml`.

**`.gitignore` is aggressive and has silently swallowed a whole directory.** It once ignored all of `packaging/`, so `git add packaging/` succeeded while committing nothing, and the Chocolatey job failed with "cannot find path" much later. After adding files under a new directory, check `git status --ignored` before assuming they're staged.

**Release asset names are a contract across four files.** `Makefile`'s `release` target produces `pull-vids-<os>-<arch>.tar.gz` and `pull-vids-windows-amd64.zip`; `install.sh`, `install.ps1`, `packaging/chocolatey/tools/chocolateyinstall.ps1`, and the README's manual-download section all hardcode those names. All four 404'd simultaneously once because they fetched bare binaries that CI never published. Rename an artifact and you must update all four.

**There is no `main` branch.** The default branch is `develop`. Raw GitHub URLs in `install.sh`/`install.ps1`/README must point at `develop` or they 404.

**`.github/copilot-instructions.md` is stale** (claims 413 lines, Go 1.18, and a `homebrew-tap/` path that isn't in this repo). Prefer this file; the formula actually lives in the separate `vib795/homebrew-tap` repo and is generated by CI.

## Release pipeline

**Pushing a `v*` tag is the only trigger.** `release.yml` builds, creates the GitHub release, then fans out to `update-homebrew.yml` and `publish-chocolatey.yml` as `workflow_call` jobs with `secrets: inherit`.

That fan-out is deliberate and must not be "simplified" back into `on: release: [published]` triggers in the child workflows. A release created with `GITHUB_TOKEN` **does not emit a `release` event** (GitHub's recursion guard), so those triggers never fire — the Homebrew tap silently went un-updated for several releases because of exactly this.

Both child workflows accept a `tag` input and derive the bare number with `${TAG#v}`, since Homebrew and Chocolatey both want the version without the `v`.

Secrets: `HOMEBREW_TAP_TOKEN` (PAT for pushing to `vib795/homebrew-tap`) and `CHOCO_API_KEY`. `publish-chocolatey.yml` deliberately **skips itself rather than failing** when `CHOCO_API_KEY` is unset, so releases stay green without it.

`packaging/chocolatey/` ships `$version$`/`$checksum$` placeholders that CI substitutes at pack time, so the repo never holds a checksum that has drifted from the current release. Chocolatey's first submission for a package ID sits in human moderation and is invisible to search until approved.

## Conventions

Release checklist: bump `version` in `main.go`, commit, tag `vX.Y.Z`, push the tag.

User-facing output goes through the package-level `color` vars (`cyan`, `green`, `yellow`, `red`, `magenta`) — never bare `fmt.Println` for messages the user reads. Errors print red and propagate to the exit code.

stdout and stderr are consumed in **separate goroutines** in `executeDownload()`; reading them serially deadlocks on large output.

`cleanURL()` strips shell-added backslash escapes (`\?`, `\=`, `\&`, `\:`, `\/`) so unquoted URLs still work.

Extension points: new quality preset → `qualityMap` in `getFormatString()` (value must be valid yt-dlp format syntax); new flag → `flag.*Var()` in `parseFlags()` plus a `Config` field, and note that every flag is registered twice for its short and long form.
