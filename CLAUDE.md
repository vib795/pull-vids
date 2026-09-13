# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Go CLI that wraps the `yt-dlp` binary: `main.go` holds the flags, download pipeline and progress parsing, and `transcript.go` holds caption handling. Both are `package main`; there are no internal packages and no abstractions — the design is intentionally flat. The program's real job is threefold: translate friendly flags into `yt-dlp` arguments, parse the subprocess's stdout into a live progress bar, and retry intelligently when a CDN pushes back.

Runtime dependencies are external binaries, not Go libraries: `yt-dlp` (checked at startup by `checkYtDlp()`), `ffmpeg` (required by yt-dlp, never checked), and `aria2c` (optional, used only with `--downloader aria2c`, never a package dependency).

## Commands

```bash
make build        # → ./pull-vids, with version injected from git describe
make build-all    # cross-compile 5 targets → dist/
make release      # wrap dist/ binaries into .tar.gz / .zip → dist/releases/
make checksums    # release + SHA256SUMS
make install      # sudo cp to /usr/local/bin
make help         # lists all targets (this is .DEFAULT_GOAL)
```

**Unit tests are narrow.** `transcript_test.go` covers caption parsing. `main_test.go` covers `retryReason`, `parseInfoLine`, and one pipeline behaviour, yt-dlp's error surviving to the retry logic, by putting a fake `yt-dlp` shell script first on `PATH`, a pattern worth reusing for other subprocess behaviour. Flags and progress parsing are untested, so a green `make test` says nothing about them. Those changes are verified by running the binary:

```bash
./pull-vids -q 720p "https://www.youtube.com/watch?v=..."
./pull-vids --downloader aria2c -N 8 "<url>"    # exercise the opt-in aria2c path
./speedtest.sh                                   # build fresh, benchmark native vs aria2c
./pull-vids -t "<url>"                            # transcript only; also try -f srt and --sub-langs xx
```

## Throughput architecture

This is the part most likely to be misunderstood. YouTube throttles by **request shape, not by connection**. A single request for a large file, whether unranged or an open-ended `Range: bytes=0-`, is held to roughly playback speed; bounded ranges of 10 MiB or less come back at full speed even on one connection. Measured with curl on one 230 MiB file: unranged 1.8 MiB/s, `bytes=0-` 1.8 MiB/s, a 10 MiB range 42.7 MiB/s. An earlier version of this file claimed Google throttled each TCP connection and that throughput scaled with connection count; that premise was wrong and led to shipping the slower backend as the default.

`resolveDownloader()` therefore maps `auto` to **native**, and aria2c is used only when named:

- **native** (default) → `--concurrent-fragments N` plus `--http-chunk-size 10M`. The chunk size is what makes YouTube fast, because it keeps every request bounded; an empty `--http-chunk-size` drops back to playback speed. `-N` only parallelises formats that are genuinely fragmented (HLS/DASH).
- **aria2c** (opt-in) → `--downloader-args aria2c:-xN -sN -jN -k1M ...`. End to end it measured 2.6–4× slower at every size tried (19 s clip: 5.1 s vs 22.3 s; 32 MiB 1080p: 6.2 s vs 19.3 s; 232 MiB 4K: 13.3 s vs 34.8 s), partly because files below its split size go out as one unranged request. Don't make it the default again, and don't add aria2 back as a Homebrew or Chocolatey dependency, without re-measuring with `./speedtest.sh`.
  - **`-j` is separate from `-x`/`-s` and critical.** For DASH/HLS formats yt-dlp hands aria2c a *fragment list*; without `-j` aria2c fetches those fragments serially. This was shipped broken in v0.3.0 and fixed in v0.3.4 — do not drop `-j` when editing these args.
  - `-k1M` is required because aria2c refuses to split ranges under 20M by default.

`downloadVideo()` wraps `executeDownload()` in a retry loop driven by `retryReason()`, which treats **HTTP 403 as a concurrency signal, not a dead URL** — it halves `config.Connections` before each retry rather than retrying identically — and retries 429s and YouTube's bot check with backoff. Its patterns are exact strings on purpose: it once matched any error containing `bot`, so a URL with "robot" in it sat through 3.5 minutes of backoff before failing. When the bot check survives every retry, the final error suggests `--cookies-from-browser`.

**Titles come from the download run itself**, via `--print video:` with the `infoLinePrefix` marker, parsed by `parseInfoLine()` in the stdout goroutine. A separate `--dump-json` call used to run first, costing ~1.7 s and a second YouTube request while ignoring `--cookies`. `--print` implies `--quiet` and `--simulate`, which would silence the progress lines and skip the download, so **`--no-quiet --no-simulate` must stay alongside it**. Both have existed since at least yt-dlp 2024.03.10, the Chocolatey dependency floor.

Because the backend is switchable, **two progress formats must be parsed** in the stdout goroutine: yt-dlp's `[download] 45.3% of 12.34MiB at 1.23MiB/s ETA 00:10` lines, and aria2c's `[#8a1b2c 12MiB/100MiB(12%) CN:16 DL:25MiB ETA:3s]` (via the `aria2Progress` regex). Adding a backend means adding a third parser.

## Transcript mode

`-t/--transcript` reuses the same `executeDownload()` pipeline with `--skip-download`, and differs in ways that each fixed or prevented a real failure:

- **Captions land in a private temp dir, then get copied to the output dir.** yt-dlp exits 0 when a video has no captions in the requested language, so an empty temp dir is the only signal for reporting "no captions found" instead of silently succeeding.
- **No `-f` quality selector and no throughput flags.** Format selection still runs under `--skip-download`, so a quality filter the video can't meet would fail a caption fetch. In this mode `-f` names the caption format (`txt`/`srt`/`vtt`), which is why `--merge-output-format` is guarded off.
- **`txt` is fetched as VTT and flattened in Go** by `captionText`. YouTube's auto-generated captions scroll, so each cue repeats the previous line. The dedup rule is deliberately narrow: a repeat is dropped only when it opens a cue starting at the *exact millisecond* the previous cue ended. Measured on real tracks, all 644 scrolling repeats in a 14-minute auto-generated track matched that, while a chorus line sung twice in uploaded captions came 760ms apart. The obvious simplification, "skip any line equal to the previous one", silently deletes those real repeats. Don't reintroduce it.
- **The "no captions" error includes a ready-to-run `yt-dlp --list-subs` command with the URL.** That's safe only because `retryReason()` matches exact phrases; loosening its patterns could make this unrecoverable error retry.
- **HTTP 429 is retried.** YouTube's caption endpoint rate-limits far more readily than video streams.

The main fixture in `transcript_test.go` is a verbatim slice of a real YouTube track, and its **whitespace-only lines are load-bearing**. Tools that strip trailing whitespace have already erased them once, so the test refuses to run without them.

## Traps

**The Makefile must build the package (`.`), not `main.go`.** `go build main.go` compiles that one file only, so it fails with undefined symbols as soon as code lives in a second file. CI runs these targets, so this breaks releases, not just local builds.

**The Makefile strips the tag's leading `v` before injecting `main.version`**, because the banner prints `v%s` itself. Before v0.5.0 releases showed `vv0.4.1`.

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

The nuspec's `iconUrl` must stay a CDN URL **pinned to a commit SHA** — Chocolatey rejects raw GitHub URLs, and a branch URL would let the bytes behind an already-published version change. `assets/icon.png` is generated by `assets/make-icon.py` (standard library only, no image tooling installed anywhere); regenerating it means committing the PNG first, then repointing `iconUrl` at that new commit.

## Conventions

Release checklist: bump `version` in `main.go`, commit, tag `vX.Y.Z`, push the tag.

User-facing output goes through the package-level `color` vars (`cyan`, `green`, `yellow`, `red`, `magenta`) — never bare `fmt.Println` for messages the user reads. Errors print red and propagate to the exit code.

stdout and stderr are consumed in **separate goroutines** in `executeDownload()`; reading them serially deadlocks on large output. Both goroutines must also reach EOF (the `readers` WaitGroup) **before** `cmd.Wait()`, which closes the pipes the moment the process exits and discards anything unread. yt-dlp prints its error last, so that is the output lost: before v0.4.1 roughly 1 in 300 failures came back as a bare `exit status 1`, and the 403/429 retry never fired.

`cleanURL()` strips shell-added backslash escapes (`\?`, `\=`, `\&`, `\:`, `\/`) so unquoted URLs still work.

Extension points: new quality preset → `qualityMap` in `getFormatString()` (value must be valid yt-dlp format syntax); new flag → `flag.*Var()` in `parseFlags()` plus a `Config` field, and note that every flag is registered twice for its short and long form.
