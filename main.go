package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// version is overridden at build time via -ldflags "-X main.version=...".
// It must stay a var: the linker cannot patch a const, so declaring it const
// silently ignores the injected tag and ships the fallback value below.
var version = "0.3.1"

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
}

// resolveDownloader decides which transfer backend to use.
// "auto" prefers aria2c when installed, because aria2c splits a single
// contiguous media URL into parallel ranged requests. yt-dlp's native
// downloader can only parallelise formats that are actually fragmented.
func resolveDownloader(choice string) string {
	switch choice {
	case "aria2c":
		return "aria2c"
	case "native":
		return "native"
	default:
		if _, err := exec.LookPath("aria2c"); err == nil {
			return "aria2c"
		}
		return "native"
	}
}

type VideoInfo struct {
	Title    string  `json:"title"`
	Duration float64 `json:"duration"`
	Uploader string  `json:"uploader"`
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

func getVideoInfo(url string) (*VideoInfo, error) {
	cmd := exec.Command("yt-dlp", "--dump-json", "--no-playlist", url)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var info VideoInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func downloadVideo(config *Config) error {
	const maxRetries = 3
	const baseWaitTime = 30 // seconds

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			waitTime := baseWaitTime * (1 << (attempt - 1)) // Exponential backoff: 30s, 60s, 120s
			yellow.Printf("\n⚠️  Rate limit detected. Waiting %d seconds before retry %d/%d...\n", waitTime, attempt, maxRetries)
			time.Sleep(time.Duration(waitTime) * time.Second)
			cyan.Println("Retrying download...")
		}

		err := executeDownload(config)
		if err == nil {
			return nil
		}

		// A 403 here is the CDN pushing back on concurrency, not a dead URL.
		// Halving the connection count on each retry usually clears it, so back
		// the parallelism off before sleeping rather than retrying identically.
		if strings.Contains(err.Error(), "HTTP Error 403") || strings.Contains(err.Error(), "Forbidden") {
			if attempt < maxRetries {
				if config.Connections > 1 {
					config.Connections /= 2
					yellow.Printf("\n⚠️  Server rejected the request (403). Reducing to %d connection(s)...\n",
						config.Connections)
				}
				continue
			}
			return fmt.Errorf("download failed after %d retries: %w", maxRetries, err)
		}

		// Check if it's a rate limiting error
		if strings.Contains(err.Error(), "Sign in to confirm you're not a bot") ||
			strings.Contains(err.Error(), "bot") {
			if attempt < maxRetries {
				continue // Retry
			}
			return fmt.Errorf("download failed after %d retries: %w", maxRetries, err)
		}

		// For other errors, don't retry
		return err
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

	// Build output template
	outputTemplate := filepath.Join(outputDir, "%(title)s.%(ext)s")
	if config.Playlist {
		outputTemplate = filepath.Join(outputDir, "%(playlist)s", "%(playlist_index)s - %(title)s.%(ext)s")
	}

	if !config.Playlist {
		cyan.Printf("Starting download from: %s\n", config.URL)
		cyan.Printf("Output directory: %s\n", outputDir)

		if config.AudioOnly {
			yellow.Println("Mode: Audio only")
		} else {
			yellow.Printf("Quality: %s\n", config.Quality)
		}
		fmt.Println()
	}

	// Get video info first
	if !config.Playlist {
		info, err := getVideoInfo(config.URL)
		if err == nil {
			mins := int(info.Duration) / 60
			secs := int(info.Duration) % 60
			magenta.Printf("Title: %s\n", info.Title)
			magenta.Printf("Duration: %dm %ds\n", mins, secs)
			fmt.Println()
		}
	}

	// Build yt-dlp command
	args := []string{
		"--newline",
		"--progress",
		"-f", getFormatString(config.Quality, config.AudioOnly),
		"-o", outputTemplate,
	}

	// Throughput tuning.
	//
	// Google's CDN rate-limits each TCP connection independently (measured at
	// roughly 3 MB/s per connection), so a single-stream download leaves a fast
	// link almost entirely idle. Aggregate throughput scales close to linearly
	// with the number of concurrent connections, which is what these flags buy.
	conns := config.Connections
	if conns < 1 {
		conns = 1
	}

	switch resolveDownloader(config.Downloader) {
	case "aria2c":
		// aria2c parallelises any single URL via ranged requests, so it speeds
		// up contiguous and fragmented formats alike. -k sets the split size;
		// without it aria2c refuses to split ranges smaller than 20M.
		args = append(args,
			"--downloader", "aria2c",
			"--downloader-args", fmt.Sprintf(
				"aria2c:-x%d -s%d -k1M --file-allocation=none --console-log-level=warn --summary-interval=1",
				conns, conns),
		)
	default:
		// Native downloader: parallelise fragments, and request the stream in
		// bounded chunks so a throttled connection is re-established rather
		// than being held at its degraded rate for the whole transfer.
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
	} else if config.Format != "" {
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
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        green.Sprint("█"),
			SaucerHead:    green.Sprint("█"),
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("%"),
	)

	var currentPercent int

	// Read output
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

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
	err = cmd.Wait()

	if err != nil {
		// Include stderr in error for retry logic
		stderrStr := stderrOutput.String()
		if stderrStr != "" {
			return fmt.Errorf("download failed: %s", stderrStr)
		}
		return fmt.Errorf("download failed: %w", err)
	}

	// Ensure bar shows 100%
	bar.Set(100)
	bar.Finish()
	fmt.Println()

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
	flag.StringVar(&config.Format, "f", "", "Output format (mp4, mkv, mp3, m4a, etc.)")
	flag.StringVar(&config.Format, "format", "", "Output format (mp4, mkv, mp3, m4a, etc.)")
	flag.StringVar(&config.Cookies, "cookies", "", "Path to cookies file (Netscape format)")
	flag.StringVar(&config.CookiesFromBrowser, "cookies-from-browser", "", "Extract cookies from browser (chrome, firefox, edge, safari, etc.)")
	flag.IntVar(&config.SleepInterval, "sleep-interval", 0, "Sleep interval in seconds between downloads (avoids rate limiting)")
	flag.IntVar(&config.Connections, "N", 8, "Number of parallel connections (higher = faster; backs off automatically on 403)")
	flag.IntVar(&config.Connections, "connections", 8, "Number of parallel connections (higher = faster; backs off automatically on 403)")
	flag.StringVar(&config.Downloader, "downloader", "auto", "Transfer backend: auto, native, or aria2c (aria2c is fastest)")
	flag.StringVar(&config.ChunkSize, "http-chunk-size", "10M", "Chunk size for the native downloader (empty disables chunking)")
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
		fmt.Fprintf(os.Stderr, "  # Download in 720p quality to specific directory\n")
		fmt.Fprintf(os.Stderr, "  pull-vids -q 720p -o ~/Videos \"https://vimeo.com/123456\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Use cookies from Chrome (fixes YouTube bot detection)\n")
		fmt.Fprintf(os.Stderr, "  pull-vids --cookies-from-browser chrome \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Use cookies from Safari\n")
		fmt.Fprintf(os.Stderr, "  pull-vids --cookies-from-browser safari \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Use cookies from file\n")
		fmt.Fprintf(os.Stderr, "  pull-vids --cookies cookies.txt \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Maximise speed (Google throttles each connection, so use several)\n")
		fmt.Fprintf(os.Stderr, "  pull-vids -N 16 \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Force the native downloader instead of aria2c\n")
		fmt.Fprintf(os.Stderr, "  pull-vids --downloader native -N 8 \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
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
