package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// version is overridden at build time via -ldflags "-X main.version=...".
// It must stay a var: the linker cannot patch a const, so declaring it const
// silently ignores the injected tag and ships the fallback value below.
var version = "0.5.0"

// infoLinePrefix marks the line yt-dlp prints, via --print, with a video's
// duration and title. Nothing yt-dlp prints on its own starts with it.
const infoLinePrefix = "pull-vids:info\t"

// botCheck is the error YouTube returns when it wants a signed-in session.
const botCheck = "Sign in to confirm you're not a bot"

// aria2Progress matches aria2c's status line, capturing percent, connection
// count, download rate and ETA:
//
//	[#8a1b2c 12MiB/100MiB(12%) CN:16 DL:25MiB ETA:3s]
var aria2Progress = regexp.MustCompile(
	`\((\d+)%\).*?CN:(\d+).*?DL:\s*([0-9.]+[KMG]?i?B)(?:.*?ETA:\s*(\S+?))?\]`)

var (
	cyan    = color.New(color.FgCyan)
	green   = color.New(color.FgGreen)
	yellow  = color.New(color.FgYellow)
	red     = color.New(color.FgRed)
	magenta = color.New(color.FgMagenta)
)

type Config struct {
	URL                string
	Output             string
	Quality            string
	AudioOnly          bool
	Playlist           bool
	Format             string
	NoBanner           bool
	ShowVersion        bool
	Cookies            string
	CookiesFromBrowser string
	SleepInterval      int
	Connections        int
	Downloader         string
	ChunkSize          string
	Transcript         bool
	Subs               bool
	SubLangs           string
}

// resolveDownloader decides which transfer backend to use. "auto" means
// native, and aria2c is used only when asked for by name.
//
// YouTube throttles requests by their shape rather than by connection: a
// single request for a large file is held to roughly playback speed, while
// bounded ranges of 10 MiB or less arrive at full speed even on one
// connection. The native downloader keeps every request bounded via
// --http-chunk-size. aria2c measured 2.6-4x slower end to end at every size
// tried, from a 19-second clip to a 232 MiB 4K video, partly because anything
// below its split size goes out as one unranged request.
func resolveDownloader(choice string) string {
	if choice == "aria2c" {
		return "aria2c"
	}
	return "native"
}

// retryReason classifies a failed download by yt-dlp's error text:
// "forbidden" for a 403, which usually means too many connections,
// "rate-limited" for an explicit 429 or YouTube's bot check, and "" for
// anything retrying won't fix.
//
// The patterns are deliberately exact. This once matched any error containing
// "bot", so a URL or title with "robot" or "bottle" in it spent three and a
// half minutes backing off from an error that could never succeed.
func retryReason(msg string) string {
	switch {
	case strings.Contains(msg, "HTTP Error 403"), strings.Contains(msg, "Forbidden"):
		return "forbidden"
	case strings.Contains(msg, "HTTP Error 429"), strings.Contains(msg, botCheck):
		return "rate-limited"
	}
	return ""
}

// parseInfoLine extracts the title and a readable duration from a line
// starting with infoLinePrefix. duration is empty when the site reports none,
// as for a live stream, where yt-dlp prints "NA".
func parseInfoLine(line string) (title, duration string, ok bool) {
	rest, found := strings.CutPrefix(line, infoLinePrefix)
	if !found {
		return "", "", false
	}
	secs, title, found := strings.Cut(rest, "\t")
	if !found {
		return "", "", false
	}
	if s, err := strconv.ParseFloat(secs, 64); err == nil {
		duration = fmt.Sprintf("%dm %ds", int(s)/60, int(s)%60)
	}
	return title, duration, true
}

func printBanner() {
	cyan.Println("╔═══════════════════════════════════════╗")
	cyan.Printf("║          pull-vids v%-8s       ║\n", version)
	cyan.Println("║  Universal Video Downloader CLI 🌐    ║")
	cyan.Println("║   YouTube • Vimeo • Twitter • More    ║")
	cyan.Println("╚═══════════════════════════════════════╝")
	fmt.Println()
}

func checkYtDlp() error {
	cmd := exec.Command("yt-dlp", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp not found. Install it with: pip install yt-dlp")
	}
	return nil
}

func getFormatString(quality string, audioOnly bool) string {
	if audioOnly {
		return "bestaudio/best"
	}

	qualityMap := map[string]string{
		"best":   "bestvideo+bestaudio/best",
		"high":   "bestvideo[height<=1080]+bestaudio/best[height<=1080]",
		"medium": "bestvideo[height<=720]+bestaudio/best[height<=720]",
		"low":    "bestvideo[height<=480]+bestaudio/best[height<=480]",
		"2160p":  "bestvideo[height<=2160]+bestaudio/best[height<=2160]",
		"1440p":  "bestvideo[height<=1440]+bestaudio/best[height<=1440]",
		"1080p":  "bestvideo[height<=1080]+bestaudio/best[height<=1080]",
		"720p":   "bestvideo[height<=720]+bestaudio/best[height<=720]",
		"480p":   "bestvideo[height<=480]+bestaudio/best[height<=480]",
		"360p":   "bestvideo[height<=360]+bestaudio/best[height<=360]",
	}

	if format, ok := qualityMap[quality]; ok {
		return format
	}
	return "bestvideo+bestaudio/best"
}

func downloadVideo(config *Config) error {
	const maxRetries = 3
	const baseWaitTime = 30 // seconds

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			waitTime := baseWaitTime * (1 << (attempt - 1)) // Exponential backoff: 30s, 60s, 120s
			yellow.Printf("\n⚠️  Waiting %d seconds before retry %d/%d...\n", waitTime, attempt, maxRetries)
			time.Sleep(time.Duration(waitTime) * time.Second)
			cyan.Println("Retrying download...")
		}

		err := executeDownload(config)
		if err == nil {
			return nil
		}

		switch retryReason(err.Error()) {
		case "forbidden":
			// A 403 here is the CDN pushing back on concurrency, not a dead URL.
			// Halving the connection count on each retry usually clears it, so
			// back the parallelism off before sleeping rather than retrying
			// identically.
			if attempt < maxRetries {
				if config.Connections > 1 {
					config.Connections /= 2
					yellow.Printf("\n⚠️  Server rejected the request (403). Reducing to %d connection(s)...\n",
						config.Connections)
				}
				continue
			}
		case "rate-limited":
			// A 429 says so outright, and caption endpoints return it far more
			// readily than media streams do.
			if attempt < maxRetries {
				continue
			}
		default:
			return err
		}

		if strings.Contains(err.Error(), botCheck) {
			return fmt.Errorf("download failed after %d retries: %w\nYouTube wants a signed-in session; retry with --cookies-from-browser <browser>", maxRetries, err)
		}
		return fmt.Errorf("download failed after %d retries: %w", maxRetries, err)
	}

	return fmt.Errorf("download failed after %d retries", maxRetries)
}

func executeDownload(config *Config) error {
	// Ensure output directory exists
	outputDir := os.ExpandEnv(config.Output)
	if strings.HasPrefix(config.Output, "~") {
		home, _ := os.UserHomeDir()
		outputDir = filepath.Join(home, config.Output[1:])
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Transcripts are fetched into a private directory rather than straight
	// into outputDir. yt-dlp exits successfully when a video has no captions
	// in the requested language, so finding nothing there is the only way to
	// report that instead of silently writing nothing. It also keeps the txt
	// conversion from touching caption files the user already has.
	fetchDir := outputDir
	if config.Transcript {
		workDir, err := os.MkdirTemp("", "pull-vids-transcript-")
		if err != nil {
			return fmt.Errorf("failed to create working directory: %w", err)
		}
		defer os.RemoveAll(workDir)
		fetchDir = workDir
	}

	// Build output template
	outputTemplate := filepath.Join(fetchDir, "%(title)s.%(ext)s")
	if config.Playlist {
		outputTemplate = filepath.Join(fetchDir, "%(playlist)s", "%(playlist_index)s - %(title)s.%(ext)s")
	}

	if !config.Playlist {
		cyan.Printf("Starting download from: %s\n", config.URL)
		cyan.Printf("Output directory: %s\n", outputDir)

		if config.Transcript {
			yellow.Printf("Mode: Transcript only (%s, languages: %s)\n", config.Format, config.SubLangs)
		} else if config.AudioOnly {
			yellow.Println("Mode: Audio only")
		} else {
			yellow.Printf("Quality: %s\n", config.Quality)
		}
		fmt.Println()
	}

	// Build yt-dlp command
	args := []string{
		"--newline",
		"--progress",
		"-o", outputTemplate,
		// The download run reports each video's title itself, rather than a
		// separate extraction beforehand, which cost ~1.7s and a second YouTube
		// request and ignored --cookies. --print would otherwise imply --quiet,
		// silencing the progress lines parsed below, and --simulate.
		"--print", "video:" + infoLinePrefix + "%(duration)s\t%(title)s",
		"--no-quiet", "--no-simulate",
	}

	if config.Transcript {
		// No -f: quality means nothing without a download, yet format selection
		// still runs under --skip-download, so a filter the video can't meet
		// would fail the fetch. --ignore-no-formats-error likewise stops a video
		// whose streams are unavailable from blocking captions that do exist.
		args = append(args, "--skip-download", "--ignore-no-formats-error")
		args = append(args, subtitleArgs(config.SubLangs, transcriptFormats[config.Format])...)
	} else {
		args = append(args, "-f", getFormatString(config.Quality, config.AudioOnly))
		if config.Subs {
			args = append(args, subtitleArgs(config.SubLangs, "srt")...)
		}
	}

	// Throughput tuning. See resolveDownloader for why native is the default:
	// the request shape, not the connection count, decides YouTube's speed.
	conns := config.Connections
	if conns < 1 {
		conns = 1
	}

	switch {
	case config.Transcript:
		// Caption files are a few kilobytes, so there is no throughput to buy,
		// and handing them to aria2c would only add a process per file.
	case resolveDownloader(config.Downloader) == "aria2c":
		// -x and -s parallelise a single URL via ranged requests, which is what
		// helps on one large contiguous file. -j is separate and just as
		// important: yt-dlp hands aria2c a fragment list for DASH/HLS formats,
		// and without -j aria2c fetches those fragments one at a time, which is
		// slower than the native downloader. -k sets the split size, because
		// aria2c will not split a range smaller than 20M by default. Files below
		// even 1M still go out as one unranged request, which YouTube throttles.
		args = append(args,
			"--downloader", "aria2c",
			"--downloader-args", fmt.Sprintf(
				"aria2c:-x%d -s%d -j%d -k1M --file-allocation=none --console-log-level=warn --summary-interval=1",
				conns, conns, conns),
		)
	default:
		// Native downloader: parallelise fragments, and request the stream in
		// bounded chunks, which YouTube serves at full speed. Without a chunk
		// size the request is unbounded and throttled to about playback speed.
		args = append(args, "--concurrent-fragments", fmt.Sprintf("%d", conns))
		if config.ChunkSize != "" {
			args = append(args, "--http-chunk-size", config.ChunkSize)
		}
	}

	// Add cookie support
	if config.CookiesFromBrowser != "" {
		args = append(args, "--cookies-from-browser", config.CookiesFromBrowser)
	} else if config.Cookies != "" {
		args = append(args, "--cookies", config.Cookies)
	}

	// Add sleep interval to avoid rate limiting
	if config.SleepInterval > 0 {
		args = append(args, "--sleep-interval", fmt.Sprintf("%d", config.SleepInterval))
	}

	if !config.Playlist {
		args = append(args, "--no-playlist")
	}

	if config.AudioOnly {
		format := "mp3"
		if config.Format != "" {
			format = config.Format
		}
		args = append(args, "-x", "--audio-format", format, "--audio-quality", "192K")
	} else if config.Format != "" && !config.Transcript {
		// In transcript mode -f names the caption format, not a container.
		args = append(args, "--merge-output-format", config.Format)
	}

	args = append(args, config.URL)

	// Execute download
	cmd := exec.Command("yt-dlp", args...)

	// Create pipe for stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Create progress bar (100% scale)
	bar := progressbar.NewOptions(100,
		progressbar.OptionSetDescription(green.Sprint("Downloading")),
		progressbar.OptionSetWidth(50),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        green.Sprint("█"),
			SaucerHead:    green.Sprint("█"),
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	var currentPercent int

	// Both pipes must be read to EOF before cmd.Wait, which closes them the
	// moment the process exits and discards anything still unread. yt-dlp
	// prints its error last and then exits, so that is exactly the output at
	// risk, and without it the retry logic can't see a 403 or 429.
	var readers sync.WaitGroup
	readers.Add(2)

	// Read output
	go func() {
		defer readers.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			if title, duration, ok := parseInfoLine(line); ok {
				// In a playlist this arrives after the previous video's bar was
				// drawn, so move off that line first.
				if currentPercent > 0 {
					fmt.Println()
				}
				magenta.Printf("Title: %s\n", title)
				if duration != "" {
					magenta.Printf("Duration: %s\n", duration)
				}
				fmt.Println()
				continue
			}

			// aria2c reports progress in its own format rather than yt-dlp's:
			//   [#8a1b2c 12MiB/100MiB(12%) CN:16 DL:25MiB ETA:3s]
			// Parse it so the bar stays live when the external downloader runs.
			if m := aria2Progress.FindStringSubmatch(line); m != nil {
				var percent int
				if _, err := fmt.Sscanf(m[1], "%d", &percent); err == nil {
					if percent > currentPercent {
						bar.Set(percent)
						currentPercent = percent
					}
				}
				desc := fmt.Sprintf("%s (Speed: %s/s, Conns: %s",
					green.Sprint("Downloading"),
					cyan.Sprint(m[3]),
					magenta.Sprint(m[2]))
				if m[4] != "" {
					desc += fmt.Sprintf(", ETA: %s", yellow.Sprint(m[4]))
				}
				bar.Describe(desc + ")")
				continue
			}

			if strings.Contains(line, "[download]") {
				// Parse percentage from yt-dlp output
				// Format: "[download]  45.3% of 12.34MiB at 1.23MiB/s ETA 00:10"
				if strings.Contains(line, "%") && strings.Contains(line, "of") {
					parts := strings.Fields(line)
					for i, part := range parts {
						if strings.HasSuffix(part, "%") {
							percentStr := strings.TrimSuffix(part, "%")
							var percent float64
							if _, err := fmt.Sscanf(percentStr, "%f", &percent); err == nil {
								newPercent := int(percent)
								if newPercent > currentPercent {
									bar.Set(newPercent)
									currentPercent = newPercent
								}

								// Extract and show speed/ETA info
								var speed, eta string
								for j := i + 1; j < len(parts); j++ {
									if parts[j] == "at" && j+1 < len(parts) {
										speed = parts[j+1]
									}
									if parts[j] == "ETA" && j+1 < len(parts) {
										eta = parts[j+1]
									}
								}

								if speed != "" && eta != "" {
									desc := fmt.Sprintf("%s (Speed: %s, ETA: %s)",
										green.Sprint("Downloading"),
										cyan.Sprint(speed),
										yellow.Sprint(eta))
									bar.Describe(desc)
								}
							}
							break
						}
					}
				} else if strings.Contains(line, "has already been downloaded") {
					yellow.Println("\n✓ Already downloaded, skipping...")
					bar.Set(100)
				} else if strings.Contains(line, "Destination") {
					// New file starting
					currentPercent = 0
					bar.Reset()
				}
			}
		}
	}()

	// Read stderr for errors and capture them
	var stderrOutput strings.Builder
	go func() {
		defer readers.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			stderrOutput.WriteString(line + "\n")
			if strings.Contains(line, "ERROR") {
				fmt.Println() // New line before error
				red.Println(line)
			}
		}
	}()

	// Wait for completion
	readers.Wait()
	err = cmd.Wait()

	if err != nil {
		// Include stderr in error for retry logic
		stderrStr := stderrOutput.String()
		if stderrStr != "" {
			return fmt.Errorf("download failed: %s", stderrStr)
		}
		return fmt.Errorf("download failed: %w", err)
	}

	// Transcripts are collected before the bar is forced to 100%, so a fetch
	// that found no captions fails without first drawing a finished bar.
	var written []string
	if config.Transcript {
		written, err = collectTranscripts(fetchDir, outputDir, config.Format)
		if err != nil {
			return fmt.Errorf("failed to save transcript: %w", err)
		}
		if len(written) == 0 {
			return fmt.Errorf("no captions found for language(s) %q; run 'yt-dlp --list-subs \"%s\"' to see which exist", config.SubLangs, config.URL)
		}
	}

	// Ensure bar shows 100%
	bar.Set(100)
	bar.Finish()
	fmt.Println()

	if config.Transcript {
		green.Println("==================================================")
		green.Printf("✓ Saved %d transcript(s):\n", len(written))
		for _, path := range written {
			green.Printf("  %s\n", path)
		}
		green.Println("==================================================")
		fmt.Println()
		return nil
	}

	green.Println("==================================================")
	green.Println("✓ Download completed successfully!")
	green.Println("==================================================")
	fmt.Println()

	return nil
}

func cleanURL(url string) string {
	// Remove backslash escapes that shells add to special characters
	// These are common when users don't quote URLs properly
	replacements := []string{
		`\?`, `?`,
		`\=`, `=`,
		`\&`, `&`,
		`\:`, `:`,
		`\/`, `/`,
	}

	cleaned := url
	for i := 0; i < len(replacements); i += 2 {
		cleaned = strings.ReplaceAll(cleaned, replacements[i], replacements[i+1])
	}

	return cleaned
}

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.Output, "o", "~/Downloads/pull-vids", "Output directory")
	flag.StringVar(&config.Output, "output", "~/Downloads/pull-vids", "Output directory")
	flag.StringVar(&config.Quality, "q", "best", "Video quality (best, high, medium, low, 2160p, 1440p, 1080p, 720p, 480p, 360p)")
	flag.StringVar(&config.Quality, "quality", "best", "Video quality (best, high, medium, low, 2160p, 1440p, 1080p, 720p, 480p, 360p)")
	flag.BoolVar(&config.AudioOnly, "a", false, "Download audio only")
	flag.BoolVar(&config.AudioOnly, "audio-only", false, "Download audio only")
	flag.BoolVar(&config.Playlist, "p", false, "Download entire playlist")
	flag.BoolVar(&config.Playlist, "playlist", false, "Download entire playlist")
	flag.StringVar(&config.Format, "f", "", "Output format (mp4, mkv, mp3, m4a, etc.; with -t: txt, srt or vtt)")
	flag.StringVar(&config.Format, "format", "", "Output format (mp4, mkv, mp3, m4a, etc.; with -t: txt, srt or vtt)")
	flag.BoolVar(&config.Transcript, "t", false, "Save the transcript only, without downloading the video (plain text unless -f says otherwise)")
	flag.BoolVar(&config.Transcript, "transcript", false, "Save the transcript only, without downloading the video (plain text unless -f says otherwise)")
	flag.BoolVar(&config.Subs, "subs", false, "Also save captions as .srt files alongside the video")
	flag.StringVar(&config.SubLangs, "sub-langs", "en", "Caption languages for -t and --subs (comma-separated; yt-dlp patterns like \"en.*\" work)")
	flag.StringVar(&config.Cookies, "cookies", "", "Path to cookies file (Netscape format)")
	flag.StringVar(&config.CookiesFromBrowser, "cookies-from-browser", "", "Extract cookies from browser (chrome, firefox, edge, safari, etc.)")
	flag.IntVar(&config.SleepInterval, "sleep-interval", 0, "Sleep interval in seconds between downloads (avoids rate limiting)")
	flag.IntVar(&config.Connections, "N", 8, "Parallel fragments for HLS/DASH formats, or connections with aria2c (halved automatically on 403)")
	flag.IntVar(&config.Connections, "connections", 8, "Parallel fragments for HLS/DASH formats, or connections with aria2c (halved automatically on 403)")
	flag.StringVar(&config.Downloader, "downloader", "auto", "Transfer backend: auto (native, fastest for YouTube), native, or aria2c")
	flag.StringVar(&config.ChunkSize, "http-chunk-size", "10M", "Chunk size for the native downloader (keep it set: YouTube throttles unchunked requests)")
	flag.BoolVar(&config.NoBanner, "no-banner", false, "Don't show the banner")
	flag.BoolVar(&config.ShowVersion, "v", false, "Show version")
	flag.BoolVar(&config.ShowVersion, "version", false, "Show version")

	flag.Usage = func() {
		printBanner()
		fmt.Fprintf(os.Stderr, "Usage: pull-vids [options] <URL>\n\n")
		fmt.Fprintf(os.Stderr, "Supports 1000+ sites: YouTube, Vimeo, Twitter, TikTok, Instagram, and more!\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Download YouTube video in best quality\n")
		fmt.Fprintf(os.Stderr, "  pull-vids \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Download from Vimeo\n")
		fmt.Fprintf(os.Stderr, "  pull-vids \"https://vimeo.com/123456789\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Download Twitter/X video\n")
		fmt.Fprintf(os.Stderr, "  pull-vids \"https://twitter.com/user/status/123456\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Download TikTok video\n")
		fmt.Fprintf(os.Stderr, "  pull-vids \"https://www.tiktok.com/@user/video/123456\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Download audio only as MP3\n")
		fmt.Fprintf(os.Stderr, "  pull-vids -a \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Save just the transcript as plain text (no video download)\n")
		fmt.Fprintf(os.Stderr, "  pull-vids -t \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Transcript with timestamps, in Spanish\n")
		fmt.Fprintf(os.Stderr, "  pull-vids -t -f srt --sub-langs es \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Download a video with its captions saved alongside\n")
		fmt.Fprintf(os.Stderr, "  pull-vids --subs \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Download in 720p quality to specific directory\n")
		fmt.Fprintf(os.Stderr, "  pull-vids -q 720p -o ~/Videos \"https://vimeo.com/123456\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Use cookies from Chrome (fixes YouTube bot detection)\n")
		fmt.Fprintf(os.Stderr, "  pull-vids --cookies-from-browser chrome \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Use cookies from Safari\n")
		fmt.Fprintf(os.Stderr, "  pull-vids --cookies-from-browser safari \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Use cookies from file\n")
		fmt.Fprintf(os.Stderr, "  pull-vids --cookies cookies.txt \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # More parallel fragments for HLS/DASH streams\n")
		fmt.Fprintf(os.Stderr, "  pull-vids -N 16 \"https://www.twitch.tv/videos/123456789\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Use aria2c instead of the native downloader\n")
		fmt.Fprintf(os.Stderr, "  pull-vids --downloader aria2c -N 8 \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Download playlist with delay to avoid rate limiting\n")
		fmt.Fprintf(os.Stderr, "  pull-vids --cookies-from-browser firefox --sleep-interval 5 -p \"https://www.youtube.com/playlist?list=PLAYLIST_ID\"\n\n")
	}

	flag.Parse()

	return config
}

func main() {
	config := parseFlags()

	if config.ShowVersion {
		fmt.Printf("pull-vids %s\n", version)
		os.Exit(0)
	}

	// Show banner unless disabled
	if !config.NoBanner {
		printBanner()
	}

	// Get URL from remaining args
	args := flag.Args()
	if len(args) == 0 {
		red.Println("✗ Error: Video URL required")
		fmt.Println("\nRun 'pull-vids -h' for usage information")
		os.Exit(1)
	}

	config.URL = cleanURL(args[0])

	// Basic URL validation
	if !strings.HasPrefix(config.URL, "http://") && !strings.HasPrefix(config.URL, "https://") {
		red.Println("✗ Error: Invalid URL (must start with http:// or https://)")
		os.Exit(1)
	}

	if config.Transcript {
		if config.AudioOnly {
			red.Println("✗ Error: -t/--transcript skips the media download, so it can't be combined with -a/--audio-only")
			os.Exit(1)
		}
		format, err := transcriptFormat(config.Format)
		if err != nil {
			red.Printf("✗ Error: %v\n", err)
			os.Exit(1)
		}
		config.Format = format
	}

	// Check if yt-dlp is installed
	if err := checkYtDlp(); err != nil {
		red.Printf("✗ Error: %v\n", err)
		os.Exit(1)
	}

	// Download the video
	start := time.Now()
	if err := downloadVideo(config); err != nil {
		red.Printf("✗ Error: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(start)
	cyan.Printf("Total time: %.2fs\n", elapsed.Seconds())
}
